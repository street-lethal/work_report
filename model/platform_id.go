package model

import (
	"encoding/json"
	"os"
)

type PlatformID struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FileToPlatformID(filePath string) (*PlatformID, error) {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var platformID PlatformID
	if err := json.Unmarshal(bin, &platformID); err != nil {
		return nil, err
	}

	return &platformID, nil
}

type PlatformRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (id PlatformID) RequestData() PlatformRequestData {
	return PlatformRequestData{
		Email:    id.Email,
		Password: id.Password,
	}
}