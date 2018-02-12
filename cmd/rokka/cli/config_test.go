package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveGetConfig(t *testing.T) {
	p := filepath.Join(os.TempDir(), "rokka-config-test")
	defer os.Remove(p)
	SetConfigPath(p)

	c := Config{
		APIKey: "asdf-ASDF-bla-123",
	}
	err := SaveConfig(c)
	if err != nil {
		t.Error(err)
	}

	nc, err := GetConfig()
	if err != nil {
		t.Error(err)
	}
	if nc.APIKey != c.APIKey {
		t.Errorf("Unexpected config API key, got: '%s', wanted: '%s'", nc.APIKey, c.APIKey)
	}
}
