package external

import (
	"testing"
	"time"
)

func TestTimeoutFromConfig(t *testing.T) {
	def := 60 * time.Second
	cases := map[string]time.Duration{
		"":      def,
		"30s":   30 * time.Second,
		"2m":    2 * time.Minute,
		"bogus": def,
		"-10s":  def,
		"0s":    def,
		" 5s ":  5 * time.Second,
	}
	for in, want := range cases {
		if got := timeoutFromConfig(in, def); got != want {
			t.Fatalf("timeoutFromConfig(%q) = %v, want %v", in, got, want)
		}
	}
}
