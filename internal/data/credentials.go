package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadCredentials() (*Credentials, error) {
	path := filepath.Join(ClaudeDir(), ".credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}
