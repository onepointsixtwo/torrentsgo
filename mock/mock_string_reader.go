package mock

import (
	"io"
)

type MockStringReader struct {
	content       []byte
	contentOffset int
}

func NewMockStringReader(content string) *MockStringReader {
	return &MockStringReader{content: []byte(content), contentOffset: 0}
}

func (reader *MockStringReader) Read(p []byte) (n int, err error) {
	contentLength := len(reader.content)
	remainingLength := contentLength - reader.contentOffset

	if remainingLength > 0 {
		maxReadBytes := len(p)
		bytesToWrite := min(remainingLength, maxReadBytes)
		for i := 0; i < bytesToWrite; i++ {
			p[i] = reader.content[i+reader.contentOffset]
		}
		reader.contentOffset = reader.contentOffset + bytesToWrite
		return bytesToWrite, nil
	} else {
		return 0, io.EOF
	}
}

func min(one, two int) int {
	if one < two {
		return one
	}
	return two
}
