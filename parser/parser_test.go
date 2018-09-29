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

	metaInfo, err := ParseMetaInfo(reader)
	if err != nil {
		t.Errorf("Unexpected error parsing meta info file %v", err)
	}

	announce := metaInfo.AnnounceUrls[0]
	if announce.String() != "http://linuxtracker.org:2710/00000000000000000000000000000000/announce" {
		t.Errorf("Unexpected announce URL found for single file torrent announce: %v", announce)
	}

	//TODO: also check the creation date, encoding, and info (length, files (zeroth only) and private)
}

func TestWithMultiFileTorrentFile(t *testing.T) {
	reader, fileErr := os.Open("../mock/multi-file.torrent")
	if fileErr != nil {
		t.Errorf("Cannot run test - failed to read file %v", fileErr)
		return
	}

	metaInfo, err := ParseMetaInfo(reader)
	if err != nil {
		t.Errorf("Unexpected error parsing meta info file %v", err)
	}

	announce := metaInfo.AnnounceUrls[0]
	if announce.String() != "http://legittorrents.info:2710/announce" {
		t.Errorf("Unexpected announce URL found for single file torrent announce: %v", announce)
	}

	//TODO: check remaining contents of meta info dictionary
}
