package chromedp

import (
	"bufio"
	"context"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Pool interface {
	Allocate(context.Context) (*Browser, error)
}

func NewPool(parent context.Context, opts ...PoolOption) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	c := &Context{}

	for _, o := range opts {
		o(&c.pool)
	}

	ctx = context.WithValue(ctx, contextKey{}, c)
	return ctx, cancel
}

type PoolOption func(*Pool)

func WithExecPool(opts ...ExecPoolOption) func(*Pool) {
	return func(p *Pool) {
		ep := &ExecPool{}
		for _, o := range opts {
			o(ep)
		}
		*p = ep
	}
}

type ExecPoolOption func(*ExecPool)

type ExecPool struct{}

func (p *ExecPool) Allocate(ctx context.Context) (*Browser, error) {
	dataDir, err := ioutil.TempDir("", "chromedp-runner")
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "chromium",
		"--no-first-run",
		"--no-default-browser-check",
		"--remote-debugging-port=0",
		"--headless",
		"--user-data-dir="+dataDir,
	)
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

	return NewBrowser(wsURL)
}
