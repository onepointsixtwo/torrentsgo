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

	announce := decoded["announce"]
	if announce != "http://linuxtracker.org:2710/00000000000000000000000000000000/announce" {
		t.Errorf("Announce key should have held url value but was '%v'", announce)
	}

	creationDate := decoded["creation date"]
	if creationDate != 1537299287 {
		t.Errorf("Expected creation date to be 1537299287 but was %v", creationDate)
	}

	encoding := decoded["encoding"]
	if encoding != "UTF-8" {
		t.Errorf("Expected encoding to be 'UTF-8' but was '%v'", encoding)
	}

	infoMapUncast := decoded["info"]
	var infoMap map[string]interface{}
	switch v := infoMapUncast.(type) {
	case map[string]interface{}:
		infoMap = v
	default:
		t.Errorf("Expected info map type to be map[string]interface{} but was %v", v)
		return
	}

	length := infoMap["length"]
	if length != 1637744640 {
		t.Errorf("Expected length to be 1637744640 but read out '%v'", length)
	}

	name := infoMap["name"]
	if name != "Reborn-OS-2018.09.17-x86_64.iso" {
		t.Errorf("Read out unexpected name from info dictionary '%v'", name)
	}

	pieceLength := infoMap["piece length"]
	if pieceLength != 1048576 {
		t.Errorf("Expected piece length to be 1048576 but was %v", pieceLength)
	}

	private := infoMap["private"]
	if private != 1 {
		t.Errorf("Expected private to be 1 but was %v", private)
	}

	pieces, ok := infoMap["pieces"].(string)
	if !ok {
		t.Errorf("Could not cast pieces to string")
	}
	lengthOfPieces := len(pieces)
	if lengthOfPieces != 31240 {
		t.Errorf("Expected length of pieces to be 31240 but was '%v'", lengthOfPieces)
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

// Basic Types Tests

func TestReadIntegerValue(t *testing.T) {
	reader := mock.NewMockStringReader("123e")

	value, err := readIntegerValue(reader)
	if err != nil {
		t.Errorf("Did not expect error when reading integer value but got %v", err)
	}
	if value != 123 {
		t.Errorf("Expected read integer value to be 123 but was %v", value)
	}
}

func TestReadStringValue(t *testing.T) {
	reader := mock.NewMockStringReader("0:teststring")

	// first character was already read off - so total is 10
	str, err := readStringValue(reader, "1")
	if err != nil {
		t.Errorf("Did not expect error when reading string value but got %v", err)
	}
	if str != "teststring" {
		t.Errorf("Expected output string to be 'teststring' but was '%v'", str)
	}
}

// Lists

func TestReadListValue(t *testing.T) {
	reader := mock.NewMockStringReader("4:spam4:eggsi10e5:tests3:twoi12ee")

	list, err := readListValue(reader)
	if err != nil {
		t.Errorf("Error encountered while reading list value")
	}

	if list[0] != "spam" {
		t.Errorf("Expected list[0] to be 'spam' but was '%v'", list[0])
	}

	if list[1] != "eggs" {
		t.Errorf("Expected list[1] to be 'eggs' but was '%v'", list[1])
	}

	if list[2] != 10 {
		t.Errorf("Expected list[2] to be 10 but was %v", list[2])
	}

	if list[3] != "tests" {
		t.Errorf("Expected list[3] to be 'tests' but was '%v'", list[3])
	}

	if list[4] != "two" {
		t.Errorf("Expected list[4] to be 'two' but was '%v'", list[4])
	}

	if list[5] != 12 {
		t.Errorf("Expected list[5] to be 12 but was %v", list[5])
	}
}
