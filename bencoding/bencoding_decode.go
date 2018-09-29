package bencoding

import (
	"fmt"
	"io"
	"strconv"
)

const (
	DICTIONARY      = "d"
	LIST            = "l"
	INTEGER         = "i"
	LENGTHDELIMETER = ":"
	END             = "e"
)

func DecodeBencoding(reader io.Reader) (map[string]interface{}, error) {
	// Since we're only supporting the outer structure being a dictionary
	// we just  check the first byte is d and then proceed to read it in as a dictionary
	firstChar, err := readLengthAsString(reader, 1)
	if err != nil {
		return nil, err
	} else if firstChar != DICTIONARY {
		return nil, fmt.Errorf("Expected bencoded structure to be dictionary with opening character 'd' but was '%v'", firstChar)
	}

	return decodeDictionary(reader)
}

// Read value

func readValue(reader io.Reader) (interface{}, error) {
	firstCharacter, err := readLengthAsString(reader, 1)
	if err != nil {
		return nil, err
	}

	switch firstCharacter {
	case INTEGER:
		return readIntegerValue(reader)
	case LIST:
		return readListValue(reader)
	case DICTIONARY:
		return decodeDictionary(reader)
	case END:
		return nil, io.EOF
	default:
		return readStringValue(reader, firstCharacter)
	}
}

// Dictionary decoding

func decodeDictionary(reader io.Reader) (map[string]interface{}, error) {
	dictionary := make(map[string]interface{})
	for {
		str, value, err := readDictionaryPair(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else if len(str) == 0 {
			break
		}
		dictionary[str] = value
	}
	return dictionary, nil
}

func readDictionaryPair(reader io.Reader) (string, interface{}, error) {
	key, err := readDictionaryKey(reader)
	if err != nil {
		return "", nil, err
	}
	value, errVal := readValue(reader)
	if errVal != nil {
		return "", nil, errVal
	}
	return key, value, nil
}

func readDictionaryKey(reader io.Reader) (string, error) {
	value, err := readValue(reader)
	if err != nil {
		return "", err
	}

	strVal, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("Failed to read dictionary key - value is not string %v", value)
	}

	return strVal, nil
}

// String decoding

func readStringValue(reader io.Reader, firstLengthCharacter string) (string, error) {
	// Read up to the ':'
	lengthString := firstLengthCharacter
	for {
		readCharacter, err := readLengthAsString(reader, 1)
		if err != nil {
			return "", err
		}

		if readCharacter == LENGTHDELIMETER {
			break
		} else {
			lengthString = lengthString + readCharacter
		}
	}

	length, err := strconv.Atoi(lengthString)
	if err != nil {
		return "", err
	}

	return readLengthAsString(reader, length)
}

// Integer decoding

func readIntegerValue(reader io.Reader) (int, error) {
	integerString := ""
	for {
		character, err := readLengthAsString(reader, 1)
		if err != nil {
			return 0, err
		}

		if character == END {
			break
		} else {
			integerString = integerString + character
		}
	}

	return strconv.Atoi(integerString)
}

// List

func readListValue(reader io.Reader) ([]interface{}, error) {
	list := make([]interface{}, 0)

	for {
		value, err := readValue(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			list = append(list, value)
		}
	}

	return list, nil
}

// Raw type helpers

func readLengthAsString(reader io.Reader, length int) (string, error) {
	b, err := readLengthAsBytes(reader, length)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func readLengthAsBytes(reader io.Reader, length int) ([]byte, error) {
	b := make([]byte, length)
	n, err := reader.Read(b)
	if err != nil {
		return nil, err
	}
	bytesRead := b[:n]
	return bytesRead, nil
}
