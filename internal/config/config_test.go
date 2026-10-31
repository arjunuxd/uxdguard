package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) (string, string) {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, ".uxdguard.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return root, path
}

func TestLoadRejectsExplicitZeroMaxFileSize(t *testing.T) {
	root, _ := writeConfig(t, "version: 1\nscan:\n  max_file_size_bytes: 0\n")
	if _, _, err := Load(root, ""); err == nil {
		t.Fatal("expected validation error for explicit zero max_file_size_bytes")
	}
}

func TestLoadPreservesExplicitEmptyFailOn(t *testing.T) {
	root, _ := writeConfig(t, "version: 1\npolicy:\n  fail_on: []\n")
	cfg, _, err := Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Policy.FailOn) != 0 {
		t.Fatalf("explicit empty fail_on was replaced: %v", cfg.Policy.FailOn)
	}
}

func TestLoadPreservesExplicitEmptyExclusions(t *testing.T) {
	root, _ := writeConfig(t, "version: 1\nscan:\n  exclusions: []\n")
	cfg, _, err := Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Scan.Exclusions) != 0 {
		t.Fatalf("explicit empty exclusions were replaced: %v", cfg.Scan.Exclusions)
	}
}

func TestLoadRejectsMissingExplicitConfig(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "missing.yaml")
	if _, _, err := Load(root, path); err == nil {
		t.Fatal("expected error for missing explicit config")
	}
}

func TestLoadRejectsConfigOutsideProjectRoot(t *testing.T) {
	root := t.TempDir()
	outsideRoot := t.TempDir()
	path := filepath.Join(outsideRoot, "config.yaml")
	if err := os.WriteFile(path, []byte("version: 1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Load(root, path); err == nil {
		t.Fatal("expected error for config outside project root")
	}
}

func TestLoadAppliesDefaultsForAbsentKeys(t *testing.T) {
	root, _ := writeConfig(t, "version: 1\n")
	cfg, _, err := Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	def := Defaults()
	if cfg.Scan.MaxFileSizeBytes != def.Scan.MaxFileSizeBytes {
		t.Fatalf("max file size = %d, want %d", cfg.Scan.MaxFileSizeBytes, def.Scan.MaxFileSizeBytes)
	}
	if len(cfg.Scan.Exclusions) != len(def.Scan.Exclusions) {
		t.Fatalf("exclusions = %v, want defaults", cfg.Scan.Exclusions)
	}
	if len(cfg.Policy.FailOn) != len(def.Policy.FailOn) {
		t.Fatalf("fail_on = %v, want defaults", cfg.Policy.FailOn)
	}
}
