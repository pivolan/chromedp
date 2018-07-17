package chromedp

import (
	"context"
)

// Executor
type Executor interface {
}

// Context
type Context struct {
	context.Context

	conn Transport

	logf func(string, ...interface{})
	errf func(string, ...interface{})
}

// NewContext creates a browser context using the parent context.
func NewContext(parent context.Context, opts ...ContextOption) (context.Context, context.CancelFunc) {
	// create root context
	ctxt, cancel := context.WithCancel(parent)

	c := &Context{Context: ctxt}

	// apply opts
	for _, o := range opts {
		o(c)
	}

	return c, cancel
}

// FromContext creates a new browser context from the provided context.
func FromContext(parent context.Context) *Context {
	return nil
}

// Run runs the tasks against the provided browser context.
func Run(ctxt context.Context, tasks Tasks) error {
	ctxt = FromContext(ctxt)
	if ctxt == nil {
		return ErrInvalidContext
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
