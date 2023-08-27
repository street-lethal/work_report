package service

import (
	"os"
)

type IOService interface {
	Input(filePath string) (string, error)
	Output(works map[int][]string, filePath string) error
}

type ioService struct {
	GenerateReportDataService
}

func NewIOService(gs GenerateReportDataService) IOService {
	return &ioService{gs}
}

func (s ioService) Input(filePath string) (string, error) {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bin), nil
}
func (s ioService) Output(works map[int][]string, filePath string) error {
	report := s.GenerateReportDataService.Generate(works)

	if err := report.ToFile(filePath); err != nil {
		return err
	}

	return nil
}
