package main

import (
	"os"
	"testing"

	"github.com/TheFinalGroup/my-own-database/cmd/dbms/persistence"
)

func TestLogCreate(t *testing.T) {
	testLogFilePath := "testlog.txt"

	defer os.Remove(testLogFilePath)

	_, err := persitence.LogCreate(testLogFilePath)
	if err != nil {
		t.Fatalf("LogCreate failed: %v", err)
	}

	if _, err := os.Stat(testLogFilePath); os.IsNotExist(err) {
		t.Errorf("Log file was not created as expected.")
	}
}

func TestLogAppend(t *testing.T) {
	testLogFilePath := "testlog.txt"

	defer os.Remove(testLogFilePath)

	logFile, err := persitence.LogCreate(testLogFilePath)
	if err != nil {
		t.Fatalf("LogCreate failed: %v", err)
	}
	defer logFile.Close()

	err = persitence.LogAppend(logFile, "Test log entry")
	if err != nil {
		t.Fatalf("LogAppend failed: %v", err)
	}

	fileContents, err := os.ReadFile(testLogFilePath)
	if err != nil {
		t.Fatalf("Error reading log file: %v", err)
	}
	expectedContents := "Test log entry\n"
	if string(fileContents) != expectedContents {
		t.Errorf("Log entry was not appended as expected.")
	}
}
