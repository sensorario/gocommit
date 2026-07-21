package commit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureConfig_CreatesFileWithDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gocommit.conf.json")

	if err := EnsureConfig(path); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.OnBeforeCommit != "" || cfg.OnAfterCommit != "" || cfg.PushAfterCommit != false {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
}

func TestEnsureConfig_PreservesExistingValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gocommit.conf.json")
	os.WriteFile(path, []byte(`{"onBeforeCommit":"echo hi","onAfterCommit":"","pushAfterCommit":true}`), 0644)

	if err := EnsureConfig(path); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.OnBeforeCommit != "echo hi" {
		t.Errorf("onBeforeCommit overwritten, got %q", cfg.OnBeforeCommit)
	}
	if !cfg.PushAfterCommit {
		t.Error("pushAfterCommit overwritten")
	}
}

func TestCheckConfig_ReportsAddedKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gocommit.conf.json")

	added, err := CheckConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(added) == 0 {
		t.Error("expected added keys for a new file, got none")
	}
}

func TestCheckConfig_NoKeysAddedWhenComplete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gocommit.conf.json")
	os.WriteFile(path, []byte(`{"onBeforeCommit":"","onAfterCommit":"","pushAfterCommit":false}`), 0644)

	added, err := CheckConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(added) != 0 {
		t.Errorf("expected no added keys, got %v", added)
	}
}
