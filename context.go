package chromedp

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
)

// Executor
type Executor interface {
	Execute(context.Context, string, json.Marshaler, json.Unmarshaler) error
}

// Context
type Context struct {
	// TODO(mvdan): use WithValue instead, for layering?
	context.Context

	browser *Browser
	conn    Transport

	logf func(string, ...interface{})
	errf func(string, ...interface{})
}

// NewContext creates a browser context using the parent context.
func NewContext(parent context.Context, opts ...ContextOption) (context.Context, context.CancelFunc) {
	// create root context
	ctx, cancel := context.WithCancel(parent)

	c := &Context{Context: ctx}

	// apply opts
	for _, o := range opts {
		o(c)
	}

	return c, cancel
}

// FromContext creates a new browser context from the provided context.
func FromContext(ctx context.Context) *Context {
	c, _ := ctx.(*Context)
	return c
}

// Run runs the tasks against the provided browser context.
func Run(ctx context.Context, tasks Tasks) error {
	c := FromContext(ctx)
	if c == nil {
		return ErrInvalidContext
	}
	if c.browser == nil {
		if err := c.startProcess(); err != nil {
			return err
		}
	}
	return nil
	var th *TargetHandler
	return tasks.Do(ctx, th)
}

func (c *Context) startProcess() error {
	dataDir, err := ioutil.TempDir("", "chromedp-runner")
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(c.Context, "chromium",
		"--no-first-run",
		"--no-default-browser-check",
		"--remote-debugging-port=0",
		"--user-data-dir="+dataDir,
	)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
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
		return err
	}
	stderr.Close()

	c.browser, err = NewBrowser(wsURL)
	if err != nil {
		return err
	}
	return nil
}

// ContextOption
type ContextOption func(*Context)

// WithURL
func WithURL(urlstr string) ContextOption {
	return func(*Context) {

	}
}
