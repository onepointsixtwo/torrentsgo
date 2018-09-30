package bencoding

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncoding(t *testing.T) {
	fileName := "../testresources/single-file.torrent"
	reader, fileErr := os.Open(fileName)
	if fileErr != nil {
		t.Errorf("Cannot run test - failed to read file %v", fileErr)
	}

	decoded, err := DecodeBencoding(reader)
	if err != nil {
		t.Errorf("Error reading bencoded data '%v'", err)
	}

	encoded, err2 := EncodeBencoding(decoded)
	if err != nil {
		t.Errorf("Failed to re-encode data: %v", err2)
	}

	originalFileContents, readFileErr := ioutil.ReadFile(fileName)
	if readFileErr != nil {
		t.Errorf("Error reading original file contents %v", readFileErr)
	}

	if !bytes.Equal(encoded, originalFileContents) {
		t.Errorf("Expected original file to equal encoded contents but did not")
	}
}
