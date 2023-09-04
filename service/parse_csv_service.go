package service

import (
	"encoding/csv"
	"os"
)

type ParseCSVService interface {
	Parse(filePath string) ([][]string, error)
}

type parseCSVService struct {
}

func NewParseCSVService() ParseCSVService {
	return &parseCSVService{}
}

func (s parseCSVService) Parse(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
