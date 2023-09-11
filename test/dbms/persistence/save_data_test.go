package main

import (
	"os"
	"testing"

	"github.com/TheFinalGroup/my-own-database/cmd/dbms/persistence"
)

func TestSaveData(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile) 

	data := []byte("Hello, World!")
	err := persitence.SaveData(tempFile, data)

	if err != nil {
		t.Errorf("SaveData failed: %v", err)
	}

	savedData, err := os.ReadFile(tempFile)
	if err != nil {
		t.Errorf("Failed to read the saved file: %v", err)
	}
	if string(savedData) != string(data) {
		t.Errorf("Saved data doesn't match expected data.")
	}
}

func createTempFile(t *testing.T) string {
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "testfile-")

	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer tempFile.Close()
	return tempFile.Name()
}
