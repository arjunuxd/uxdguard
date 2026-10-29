package scan

// NoopScanner is a placeholder scanner (does nothing). Registered by default for extensibility.
type NoopScanner struct{}

func init() {
	_ = RegisterDefault(&NoopScanner{})
}

func (n *NoopScanner) ID() string          { return "noop" }
func (n *NoopScanner) Name() string        { return "No-op Scanner" }
func (n *NoopScanner) Description() string { return "No-op scanner" }
func (n *NoopScanner) Run(ctx *Context) ([]Finding, []Error) {
	return nil, nil
}
