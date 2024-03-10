package services

import (
	"strings"
	"testing"
)

func TestCreateAndRunDockerContainer(t *testing.T) {
	command := "echo Hello, World!"
	output, err := CreateAndRunDockerContainer(command)
	if err != nil {
		t.Errorf("Error running Docker container: %v", err)
	}

	expectedOutput := "Hello, World!"
	if strings.TrimSpace(output) != expectedOutput {
		t.Errorf("Expected output to be %s, got %s", expectedOutput, output)
	}
}
