package bencoding

import (
	"bytes"
	"fmt"
	"github.com/onepointsixtwo/torrentsgo/model"
	"strconv"
)

func EncodeBencoding(m *model.OrderedMap) ([]byte, error) {
	return encodeMap(m)
}

func encodeValue(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case *model.OrderedMap:
		return encodeMap(v)
	case []interface{}:
		return encodeList(v)
	case string:
		return encodeString(v)
	case int:
		return encodeInteger(v)
	default:
		return nil, fmt.Errorf("Unknown type for value %v", value)
	}
}

// TODO: need to write start and end characters for dict and list

func encodeMap(m *model.OrderedMap) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(DICTIONARY)

	var err error
	m.Iterate(func(key string, value interface{}) {
		if err != nil {
			return
		}

		keyData, keyEncodeErr := encodeString(key)
		if keyEncodeErr != nil {
			err = keyEncodeErr
			return
		}

		_, writeErr1 := buffer.Write(keyData)
		if writeErr1 != nil {
			err = writeErr1
			return
		}

		valueData, valueEncodeErr := encodeValue(value)
		if valueEncodeErr != nil {
			err = valueEncodeErr
			return
		}

		_, writeErr2 := buffer.Write(valueData)
		if writeErr2 != nil {
			err = writeErr2
			return
		}
	})

	if err != nil {
		return nil, err
	}

	buffer.WriteString(END)
	return buffer.Bytes(), nil
}

func encodeList(list []interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(LIST)

	for i := 0; i < len(list); i++ {
		value := list[i]

		encodedValue, encodingErr := encodeValue(value)
		if encodingErr != nil {
			return nil, encodingErr
		}

		_, writeErr := buffer.Write(encodedValue)
		if writeErr != nil {
			return nil, writeErr
		}
	}
	buffer.WriteString(END)
	return buffer.Bytes(), nil
}

func encodeString(value string) ([]byte, error) {
	length := len(value)
	lengthString := strconv.Itoa(length)
	return []byte(lengthString + LENGTHDELIMETER + value), nil
}

func encodeInteger(value int) ([]byte, error) {
	return []byte(INTEGER + strconv.Itoa(value) + END), nil
}
