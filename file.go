package utils

import (
	"encoding/csv"
	"os"
)

// UnpackSingleColCSV reads a single column CSV and returns the contents minus the first line
func UnpackSingleColCSV(fileName string) ([]string, error) {
	lines, err := readFromFile(fileName)
	if err != nil {
		return nil, err
	}
	var entries []string
	for i, entry := range lines {
		if i == 0 {
			continue
		}
		entries = append(entries, entry[0])
	}

	return entries, nil
}

func readFromFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
