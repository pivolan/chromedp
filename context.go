package chromedp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// Executor
type Executor interface {
	Execute(context.Context, string, json.Marshaler, json.Unmarshaler) error
}

// Context
type Context struct {
	withURL string

	pool    Pool
	browser *Browser
	handler *TargetHandler

	logf func(string, ...interface{})
	errf func(string, ...interface{})
}

// NewContext creates a browser context using the parent context.
func NewContext(parent context.Context, opts ...ContextOption) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	c := &Context{}
	if pc := FromContext(parent); pc != nil {
		c.pool = pc.pool
	}

	for _, o := range opts {
		o(c)
	}
	if c.pool == nil {
		WithExecPool()(&c.pool)
	}

	ctx = context.WithValue(ctx, contextKey{}, c)
	return ctx, cancel
}

type contextKey struct{}

// FromContext creates a new browser context from the provided context.
func FromContext(ctx context.Context) *Context {
	c, _ := ctx.Value(contextKey{}).(*Context)
	return c
}

// Run runs the action against the provided browser context.
func Run(ctx context.Context, action Action) error {
	c := FromContext(ctx)
	if c == nil || c.pool == nil {
		return ErrInvalidContext
	}
	if c.browser == nil {
		browser, err := c.pool.Allocate(ctx)
		if err != nil {
			return err
		}
		c.browser = browser
	}
	if c.handler == nil {
		if err := c.newHandler(ctx); err != nil {
			return err
		}
	}
	return action.Do(ctx, c.handler)
}

func (c *Context) newHandler(ctx context.Context) error {
	// TODO: add RemoteAddr() to the Transport interface?
	conn := c.browser.conn.(*Conn).Conn
	addr := conn.RemoteAddr()
	url := "http://" + addr.String() + "/json/new?" + url.QueryEscape(c.withURL)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var wurl withWebsocketURL
	if err := json.NewDecoder(resp.Body).Decode(&wurl); err != nil {
		return err
	}
	c.handler, err = NewTargetHandler(wurl.WebsocketURL)
	if err != nil {
		return err
	}
	if err := c.handler.Run(ctx); err != nil {
		return err
	}
	return nil
}

type withWebsocketURL struct {
	WebsocketURL string `json:"webSocketDebuggerUrl"`
}

// ContextOption
type ContextOption func(*Context)

// WithURL
func WithURL(urlstr string) ContextOption {
	return func(c *Context) { c.withURL = urlstr }
}
