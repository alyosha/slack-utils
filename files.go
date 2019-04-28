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

// CreateAndWriteCSV creates and then writes to a new CSV file
func CreateAndWriteCSV(fileName string, entries [][]string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)
	err = w.WriteAll(entries)

	if err != nil {
		return err
	}

	return nil
}

func readFromFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
