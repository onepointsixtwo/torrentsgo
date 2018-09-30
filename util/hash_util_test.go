package util

import (
	"encoding/hex"
	"testing"
)

func TestBadHashUrlEncode(t *testing.T) {
	bytes := make([]byte, 0)
	_, err := UrlEncodeHash(bytes)
	if err == nil {
		t.Errorf("Expected error while encoding hash bytes")
	}
}

func TestUrlEncodeHash(t *testing.T) {
	bytes, err := hex.DecodeString("7cd350e5a70f0a61593e636543f9fc670ffa8a4d")
	if err != nil {
		t.Errorf("Unable to decode string from hex")
	}

	encodedHash, err2 := UrlEncodeHash(bytes)
	if err2 != nil {
		t.Errorf("Unexpected error encoding hash %v", err2)
	}

	if encodedHash != "%7c%d3P%e5%a7%0f%0aaY%3eceC%f9%fcg%0f%fa%8aM" {
		t.Errorf("Unexpected value of encoded hash - '%v'", encodedHash)
	}
}
