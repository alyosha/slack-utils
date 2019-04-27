package utils

import (
	"encoding/csv"
	"errors"
	"os"
	"testing"
)

const testFileName = "test-file.csv"

var testEntries = [][]string{
	{"emails"},
	{"steve@test.com"},
	{"alyosha@test.com"},
}

func setupTestFile() error {
	_, err := os.Stat(testFileName)

	if os.IsNotExist(err) {
		file, err := os.Create(testFileName)
		if err != nil {
			return err
		}
		defer file.Close()

		w := csv.NewWriter(file)
		w.WriteAll(testEntries)

		return nil
	}

	return errors.New("file already exists")
}

func deleteTestFile() error {
	err := os.Remove(testFileName)
	if err != nil {
		return err
	}
	return nil
}

func TestUnpackSingleColCSV(t *testing.T) {
	if err := setupTestFile(); err != nil {
		t.Fatalf("failed to create test CSV file: %v", err)
	}
	defer deleteTestFile()

	emails, err := UnpackSingleColCSV(testFileName)
	if err != nil {
		t.Fatalf("failed to unpack single column CSV")
		return
	}

	if len(emails) != len(testEntries)-1 {
		t.Fatalf("expected CSV entry of length: %v, got length: %v", len(testEntries), len(emails))
		return
	}

	for i, email := range emails {
		if email != testEntries[i+1][0] {
			t.Fatalf("expected entry: %v, got %v instead", testEntries[i], email)
			return
		}
	}
}

func TestCreateAndWriteCSV(t *testing.T) {
	if err := CreateAndWriteCSV(testFileName, testEntries); err != nil {
		t.Fatalf("failed to create and write test CSV file: %v", err)
	}
	defer deleteTestFile()

	emails, err := UnpackSingleColCSV(testFileName)
	if err != nil {
		t.Fatalf("failed to unpack single column CSV")
		return
	}

	if len(emails) != len(testEntries)-1 {
		t.Fatalf("expected CSV entry of length: %v, got length: %v", len(testEntries), len(emails))
		return
	}

	for i, email := range emails {
		if email != testEntries[i+1][0] {
			t.Fatalf("expected entry: %v, got %v instead", testEntries[i], email)
			return
		}
	}
}
