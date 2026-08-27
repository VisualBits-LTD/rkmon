package collect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFirstExistingPath(t *testing.T) {
	root := t.TempDir()
	panthor := filepath.Join(root, "fb000000.gpu-panthor")
	mali := filepath.Join(root, "fb000000.gpu-mali")
	for _, path := range []string{panthor, mali} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	if got := firstExistingPath(panthor, mali); got != panthor {
		t.Fatalf("want preferred path %q, got %q", panthor, got)
	}
	if err := os.Remove(panthor); err != nil {
		t.Fatal(err)
	}
	if got := firstExistingPath(panthor, mali); got != mali {
		t.Fatalf("want fallback path %q, got %q", mali, got)
	}
}

func TestReadDevfreqUsesCurFreq(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "load"), []byte("9@200000000Hz\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cur_freq"), []byte("300000000\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	collector := New()
	t.Cleanup(collector.Close)
	got := collector.readDevfreq(root, "GPU", "gpu-thermal")
	if got.PctUsed != 9 || got.FreqHz != 300000000 {
		t.Fatalf("want load 9 and cur_freq 300000000, got %+v", got)
	}
}
