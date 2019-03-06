package chromedp

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Pool interface {
	// Allocate creates a new browser from the pool. It can be cancelled via
	// the provided context, at which point all the resources used by the
	// browser (such as temporary directories) will be cleaned up.
	Allocate(context.Context) (*Browser, error)

	// Wait can be called after cancelling a pool's context, to block until
	// all the pool's resources have been cleaned up.
	Wait()
}

func NewPool(parent context.Context, opts ...PoolOption) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	c := &Context{}

	for _, o := range opts {
		o(&c.Pool)
	}

	ctx = context.WithValue(ctx, contextKey{}, c)
	return ctx, cancel
}

type PoolOption func(*Pool)

func WithExecPool(opts ...ExecPoolOption) func(*Pool) {
	return func(p *Pool) {
		ep := &ExecPool{
			initFlags: make(map[string]interface{}),
		}
		for _, o := range opts {
			o(ep)
		}
		*p = ep
	}
}

type ExecPoolOption func(*ExecPool)

type ExecPool struct {
	initFlags map[string]interface{}

	dataDirs map[*Browser]string
	wg       sync.WaitGroup
}

func (p *ExecPool) Allocate(ctx context.Context) (*Browser, error) {
	removeDir := false
	var cmd *exec.Cmd

	// TODO: figure out a nicer way to do this
	flags := make(map[string]interface{})
	for name, value := range p.initFlags {
		flags[name] = value
	}

	dataDir, ok := flags["user-data-dir"].(string)
	if !ok {
		tempDir, err := ioutil.TempDir("", "chromedp-runner")
		if err != nil {
			return nil, err
		}
		flags["user-data-dir"] = tempDir
		dataDir = tempDir
		removeDir = true
	}

	p.wg.Add(1)
	go func() {
		<-ctx.Done()
		// First wait for the process to be finished.
		if cmd != nil {
			cmd.Wait()
		}
		// Then delete the temporary user data directory, if needed.
		if removeDir {
			os.RemoveAll(dataDir)
		}
		p.wg.Done()
	}()

	flags["remote-debugging-port"] = "0"

	prog := "chromium"
	args := []string{}
	for name, value := range flags {
		switch value := value.(type) {
		case string:
			args = append(args, fmt.Sprintf("--%s=%s", name, value))
		case bool:
			if value {
				args = append(args, fmt.Sprintf("--%s", name))
			}
		default:
			return nil, fmt.Errorf("invalid exec pool flag")
		}
	}

	cmd = exec.CommandContext(ctx, prog, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Pick up the browser's websocket URL from stderr.
	wsURL := ""
	scanner := bufio.NewScanner(stderr)
	prefix := "DevTools listening on"
	for scanner.Scan() {
		line := scanner.Text()
		if s := strings.TrimPrefix(line, prefix); s != line {
			wsURL = strings.TrimSpace(s)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	stderr.Close()

	browser, err := NewBrowser(wsURL)
	if err != nil {
		return nil, err
	}
	browser.UserDataDir = dataDir
	return browser, nil
}

func (p *ExecPool) Wait() {
	p.wg.Wait()
}

// Flag is a generic command line option to pass a flag to Chrome. If the value
// is a string, it will be passed as --name=value. If it's a boolean, it will be
// passed as --name if value is true.
func Flag(name string, value interface{}) ExecPoolOption {
	return func(p *ExecPool) {
		p.initFlags[name] = value
	}
}

// UserDataDir is the command line option to set the user data dir.
//
// Note: set this option to manually set the profile directory used by Chrome.
// When this is not set, then a default path will be created in the /tmp
// directory.
func UserDataDir(dir string) ExecPoolOption {
	return Flag("user-data-dir", dir)
}

// ProxyServer is the command line option to set the outbound proxy server.
func ProxyServer(proxy string) ExecPoolOption {
	return Flag("proxy-server", proxy)
}

// WindowSize is the command line option to set the initial window size.
func WindowSize(width, height int) ExecPoolOption {
	return Flag("window-size", fmt.Sprintf("%d,%d", width, height))
}

// UserAgent is the command line option to set the default User-Agent
// header.
func UserAgent(userAgent string) ExecPoolOption {
	return Flag("user-agent", userAgent)
}

// NoSandbox is the Chrome comamnd line option to disable the sandbox.
func NoSandbox(p *ExecPool) {
	Flag("no-sandbox", true)(p)
}

// NoFirstRun is the Chrome comamnd line option to disable the first run
// dialog.
func NoFirstRun(p *ExecPool) {
	Flag("no-first-run", true)(p)
}

// NoDefaultBrowserCheck is the Chrome comamnd line option to disable the
// default browser check.
func NoDefaultBrowserCheck(p *ExecPool) {
	Flag("no-default-browser-check", true)(p)
}

// Headless is the command line option to run in headless mode.
func Headless(p *ExecPool) {
	Flag("headless", true)(p)
}

// DisableGPU is the command line option to disable the GPU process.
func DisableGPU(p *ExecPool) {
	Flag("disable-gpu", true)(p)
}
