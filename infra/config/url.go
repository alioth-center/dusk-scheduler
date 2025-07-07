package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
)

type remoteURLConfig struct {
	client *http.Client
}

func NewRemoteURLConfig() Config {
	return &remoteURLConfig{client: http.DefaultClient}
}

func (rc *remoteURLConfig) ParseAppConfig(source, _ string, receiver any) error {
	response, executeErr := rc.client.Get(source)
	if executeErr != nil {
		return executeErr
	}

	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "application/json"):
		return json.NewDecoder(response.Body).Decode(receiver)
	case strings.Contains(contentType, "text/yaml"):
		return yaml.NewDecoder(response.Body).Decode(receiver)
	default:
		return fmt.Errorf("unexpected content type: %s", contentType)
	}
}
