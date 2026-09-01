package commit

import (
	"encoding/json"
	"os"
)

type Config struct {
	OnBeforeCommit  string `json:"onBeforeCommit"`
	OnAfterCommit   string `json:"onAfterCommit"`
	PushAfterCommit bool   `json:"pushAfterCommit"`
	MainBranch      string `json:"mainBranch,omitempty"`
}

// configDefaults maps each JSON key to its default value.
var configDefaults = map[string]interface{}{
	"onBeforeCommit":  "",
	"onAfterCommit":   "",
	"pushAfterCommit": false,
}

func EnsureConfig(path string) error {
	raw := make(map[string]json.RawMessage)
	needsWrite := false

	if _, err := os.Stat(path); os.IsNotExist(err) {
		needsWrite = true
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
	}

	for key, defaultVal := range configDefaults {
		if _, exists := raw[key]; !exists {
			encoded, err := json.Marshal(defaultVal)
			if err != nil {
				return err
			}
			raw[key] = json.RawMessage(encoded)
			needsWrite = true
		}
	}

	if needsWrite {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(raw); err != nil {
			return err
		}
	}

	return nil
}

// CheckConfig ensures the config file exists and all variables are present.
// It returns the list of keys that were added with their default values.
func CheckConfig(path string) ([]string, error) {
	raw := make(map[string]json.RawMessage)
	var added []string
	needsWrite := false

	if _, err := os.Stat(path); os.IsNotExist(err) {
		needsWrite = true
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
	}

	for key, defaultVal := range configDefaults {
		if _, exists := raw[key]; !exists {
			encoded, err := json.Marshal(defaultVal)
			if err != nil {
				return nil, err
			}
			raw[key] = json.RawMessage(encoded)
			added = append(added, key)
			needsWrite = true
		}
	}

	if needsWrite {
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(raw); err != nil {
			return nil, err
		}
	}

	return added, nil
}

func SaveConfig(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
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
