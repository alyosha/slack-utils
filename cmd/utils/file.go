package utils

import (
	"encoding/csv"
	"log"
	"os"
)

// UnpackSingleColCSV reads a single column CSV and returns the contents minus the first line
func UnpackSingleColCSV(fileName string) []string {
	lines := readFromFile(fileName)
	var entries []string
	for i, entry := range lines {
		if i == 0 {
			continue
		}
		entries = append(entries, entry[0])
	}

	return entries
}

func readFromFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Printf("error opening csv file %v", err)
		return nil
	}

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Printf("error reading csv file %v", err)
		return nil
	}

	return lines
}
