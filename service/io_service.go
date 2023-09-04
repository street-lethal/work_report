package service

import (
	"os"
)

type IOService interface {
	Input(filePath string) (string, error)
}

type ioService struct {
}

func NewIOService() IOService {
	return &ioService{}
}

func (s ioService) Input(filePath string) (string, error) {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bin), nil
}
