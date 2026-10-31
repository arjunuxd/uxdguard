package root

import "testing"

func TestNormalizeVersion(t *testing.T) {
	cases := map[string]string{
		"":          "",
		"v1.1.0":    "1.1.0",
		"1.1.0":     "1.1.0",
		"  v1.2.0 ": "1.2.0",
	}
	for in, want := range cases {
		if got := normalizeVersion(in); got != want {
			t.Fatalf("normalizeVersion(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLinkerMetadataOverridesBuildInfo(t *testing.T) {
	SetVersion("v1.1.0")
	SetBuildMetadata("abc123", "2026-10-31T00:00:00Z")
	info := ResolveBuildInfo()
	if info.Version != "1.1.0" {
		t.Fatalf("version = %q, want 1.1.0", info.Version)
	}
	if info.Commit != "abc123" {
		t.Fatalf("commit = %q, want abc123", info.Commit)
	}
	if info.BuildDate != "2026-10-31T00:00:00Z" {
		t.Fatalf("build date = %q", info.BuildDate)
	}
}
