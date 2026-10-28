package scan

import (
	"context"

	"github.com/arjunuxd/uxdguard/internal/config"
)

// Context provides data to scanners.
type Context struct {
	Root  string
	Files []string
	Cfg   *config.Config
	Ctx   context.Context
}

// NewContext creates a new scan context.
func NewContext(root string, cfg *config.Config, files []string) *Context {
	return &Context{
		Root:  root,
		Files: files,
		Cfg:   cfg,
		Ctx:   context.Background(),
	}
}

// WithContext sets the context.
func (c *Context) WithContext(ctx context.Context) *Context {
	c.Ctx = ctx
	return c
}
