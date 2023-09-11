package dbms_test

import (
	"testing"

	"github.com/TheFinalGroup/my-own-database/cmd/dbms"
)

func TestAppGreeting(t *testing.T) {
	expectedOutput := "Welcome to My Own Database Project!"
	actualOutput := dbms.AppGreeting()

	if actualOutput != expectedOutput {
		t.Errorf("Expected output: %q, but got: %q", expectedOutput, actualOutput)
	}
}
