package root

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// isTerminal reports whether w is attached to a character device.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// progress renders an in-place status line on interactive terminals. On
// non-interactive writers it stays silent so logs and pipelines remain clean.
type progress struct {
	out         io.Writer
	interactive bool
	wrote       bool
	width       int
}

func newProgress(out io.Writer, enabled bool) *progress {
	return &progress{out: out, interactive: enabled && isTerminal(out)}
}

// Status updates the in-place status line.
func (p *progress) Status(format string, args ...any) {
	if !p.interactive {
		return
	}
	line := fmt.Sprintf(format, args...)
	if len(line) > p.width {
		p.width = len(line)
	}
	padding := strings.Repeat(" ", p.width-len(line))
	fmt.Fprintf(p.out, "\r%s%s", line, padding)
	p.wrote = true
}

// End clears the status line so normal output starts on a clean row.
func (p *progress) End() {
	if !p.wrote {
		return
	}
	fmt.Fprintf(p.out, "\r%s\r", strings.Repeat(" ", p.width))
	p.wrote = false
	p.width = 0
}
