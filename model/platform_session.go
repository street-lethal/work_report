package model

import (
	"encoding/json"
	"os"
)

type PlatformSession struct {
	SessionID string `json:"session_id"`
	AWSAuth   string `json:"aws_auth"`
}

func FileToPlatformSession(filePath string) (*PlatformSession, error) {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var session PlatformSession
	if err := json.Unmarshal(bin, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s PlatformSession) ToFile(filePath string) error {
	file, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, file, 0644); err != nil {
		return err
	}

	return nil
}

func (s PlatformSession) Empty() bool {
	return s.SessionID == "" || s.AWSAuth == ""
}
