package util

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

const (
	EXPECTED_HASH_LENGTH = 20
)

func UrlEncodeHash(hashedData []byte) (string, error) {
	// Encode URI component golang equivalent is net/util's QueryEscape(string)string
	// Bytes len should be 20 if the hash is valid
	hashLen := len(hashedData)
	if hashLen != EXPECTED_HASH_LENGTH {
		return "", fmt.Errorf("Valid hash must have a byte length of %v but was %v", EXPECTED_HASH_LENGTH, hashLen)
	}

	output := ""
	for i := 0; i < hashLen; i++ {
		b := hashedData[i]

		strVal := ""
		if b <= 127 {
			stringValueOfByte := fmt.Sprintf("%c", b)
			strVal = url.QueryEscape(stringValueOfByte)
			if strings.HasPrefix(strVal, "%") {
				strVal = strings.ToLower(strVal)
			}
		} else {
			strVal = "%" + byteToHexString(b)
		}

		output = output + strVal
	}

	return output, nil
}

func byteToHexString(b byte) string {
	singleByteSlice := make([]byte, 0)
	singleByteSlice = append(singleByteSlice, b)

	outputSliceLen := hex.EncodedLen(len(singleByteSlice))
	outputSlice := make([]byte, outputSliceLen)
	hex.Encode(outputSlice, singleByteSlice)

	return string(outputSlice)
}
