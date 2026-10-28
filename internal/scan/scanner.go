package scan

// Scanner is the interface all scanners implement.
type Scanner interface {
	ID() string
	Name() string
	Description() string
	Run(ctx *Context) ([]Finding, []Error)
}

// Error represents an operational error during scanning.
type Error struct {
	ScannerID string `json:"scanner_id"`
	Message   string `json:"message"`
}
