package scan

import (
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
)

func TestScannerMaxTokenSizeDefaultsToMinimum(t *testing.T) {
	if got := scannerMaxTokenSize(nil); got != 64*1024 {
		t.Fatalf("nil config = %d, want 64KiB", got)
	}
	cfg := config.Defaults()
	cfg.Scan.MaxFileSizeBytes = 1024
	if got := scannerMaxTokenSize(cfg); got != 64*1024 {
		t.Fatalf("small max file size = %d, want 64KiB", got)
	}
}

func TestScannerMaxTokenSizeTracksFileLimit(t *testing.T) {
	cfg := config.Defaults()
	cfg.Scan.MaxFileSizeBytes = 2 * 1024 * 1024
	if got := scannerMaxTokenSize(cfg); got != 2*1024*1024 {
		t.Fatalf("max token size = %d, want file limit", got)
	}
}

func TestScannerMaxTokenSizeIsCapped(t *testing.T) {
	cfg := config.Defaults()
	cfg.Scan.MaxFileSizeBytes = 1024 * 1024 * 1024
	if got := scannerMaxTokenSize(cfg); got != 16*1024*1024 {
		t.Fatalf("max token size = %d, want 16MiB cap", got)
	}
}
