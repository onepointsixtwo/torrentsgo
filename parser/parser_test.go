package parser

import (
	"os"
	"testing"
)

func TestWithSingleFileTorrentFile(t *testing.T) {
	reader, fileErr := os.Open("../mock/linux_test_torrent.torrent")
	if fileErr != nil {
		t.Errorf("Cannot run test - failed to read file %v", fileErr)
		return
	}

	_, err := ParseMetaInfo(reader)
	if err != nil {
		t.Errorf("Unexpected error parsing meta info file %v", err)
	}

}

func TestWithMultiFileTorrentFile(t *testing.T) {
	reader, fileErr := os.Open("../mock/multi-file.torrent")
	if fileErr != nil {
		t.Errorf("Cannot run test - failed to read file %v", fileErr)
		return
	}

	_, err := ParseMetaInfo(reader)
	if err != nil {
		t.Errorf("Unexpected error parsing meta info file %v", err)
	}
}
