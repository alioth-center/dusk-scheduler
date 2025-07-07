package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type fileConfig struct {
	filename string
}

func NewFileConfig(filename string) Config {
	return &fileConfig{filename: filename}
}

func (f *fileConfig) ParseAppConfig(_, _ string, receiver any) error {
	file, readErr := os.ReadFile(f.filename)
	if readErr != nil {
		return readErr
	}

	switch filepath.Ext(f.filename) {
	case ".yaml", ".yml":
		return yaml.Unmarshal(file, receiver)
	case ".json":
		return json.Unmarshal(file, receiver)
	default:
		panic("unknown config file extension")
	}

}
