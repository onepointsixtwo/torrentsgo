package mock

import (
	"bytes"
	"io"
	"testing"
)

func TestMockStringReader(t *testing.T) {
	var reader io.Reader
	reader = NewMockStringReader("MockString")

	// Read part one
	partOneBytes := make([]byte, 4)
	n, err := reader.Read(partOneBytes)

	if err != nil || n != 4 || !bytes.Equal([]byte("Mock"), partOneBytes) {
		t.Errorf("Expected output to be string 'Mock' but was %v", string(partOneBytes))
	}

	// Read part two
	partTwoBytes := make([]byte, 6)
	n, err = reader.Read(partTwoBytes)

	if err != nil || n != 6 || !bytes.Equal([]byte("String"), partTwoBytes) {
		t.Errorf("Expected output to be string 'String' but was %v", string(partTwoBytes))
	}

	// Check for EOF
	_, err = reader.Read(make([]byte, 1))
	if err != io.EOF {
		t.Error("Expected io.EOF after reading all data from mock string reader")
	}
}
