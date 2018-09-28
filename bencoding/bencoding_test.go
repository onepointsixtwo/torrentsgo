package bencoding

import (
	"bytes"
	"github.com/onepointsixtwo/torrentsgo/mock"
	"os"
	"testing"
)

// Entire file read test

func TestReadingValidBencodedFile(t *testing.T) {
	reader, fileErr := os.Open("./linux_test_torrent.torrent")
	if fileErr != nil {
		t.Errorf("Cannot run test - failed to read file %v", fileErr)
	}

	decoded, err := DecodeBencoding(reader)
	if err != nil {
		t.Errorf("Error reading bencoded data '%v'", err)
	}
}

// Helper tests

func TestReadLengthAsString(t *testing.T) {
	reader := mock.NewMockStringReader("TEST")

	str, err := readLengthAsString(reader, 4)
	if err != nil {
		t.Error("Should not have produced error while reading length as string")
	}
	if str != "TEST" {
		t.Errorf("Expected read string to be 'TEST' but was '%v'", str)
	}
}

func TestReadLengthAsBytes(t *testing.T) {
	reader := mock.NewMockStringReader("BYTES")

	b, err := readLengthAsBytes(reader, 5)
	if err != nil {
		t.Error("Should not have produced error while reading length as bytes")
	}
	if !bytes.Equal([]byte("BYTES"), b) {
		t.Errorf("Expected read bytes to be 'BYTES' but was '%v'", string(b))
	}
}

// Dictionary tests

func TestReadDictionaryKey(t *testing.T) {
	reader := mock.NewMockStringReader("3:KEY")

	key, err := readDictionaryKey(reader)
	if err != nil {
		t.Error("No error should occur reading dictionary key")
	}
	if key != "KEY" {
		t.Errorf("Dictionary key should have been 'KEY' but was '%v'", key)
	}
}

func TestReadDictionaryPair(t *testing.T) {
	simpleDictionaryPair := "8:announce70:http://linuxtracker.org:2710/00000000000000000000000000000000/announce"
	simpleDictionaryPairReader := mock.NewMockStringReader(simpleDictionaryPair)

	key, value, err := readDictionaryPair(simpleDictionaryPairReader)
	if err != nil {
		t.Error("Unexpected error while reading dictionary pair")
	}
	if key != "announce" {
		t.Errorf("Expected dictionary key to be 'announce' but was '%v'", key)
	}
	if value != "http://linuxtracker.org:2710/00000000000000000000000000000000/announce" {
		t.Errorf("Unexpected value for dictionary value '%v'", value)
	}
}
