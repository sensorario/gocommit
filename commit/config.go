package commit

import (
	"encoding/json"
	"os"
)

type Config struct {
	OnBeforeCommit string `json:"onBeforeCommit"`
	OnAfterCommit  string `json:"onAfterCommit"`
}

func EnsureConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		config := Config{"", ""}
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(config); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
