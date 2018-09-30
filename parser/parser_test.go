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

	// Check basic metainfo

	announce := metaInfo.AnnounceUrls[0]
	if announce.String() != "http://linuxtracker.org:2710/00000000000000000000000000000000/announce" {
		t.Errorf("Unexpected announce URL found for single file torrent announce: %v", announce)
	}

	creationDate := metaInfo.CreationDate.Unix()
	if creationDate != 1537299287 {
		t.Errorf("Expected creation date to be 1537299287 but was %v", creationDate)
	}

	encoding := metaInfo.Encoding
	if encoding != "UTF-8" {
		t.Errorf("Expected encoding to be UTF-8 but was %v", encoding)
	}

	// Check info

	info := metaInfo.Info

	pieceLength := info.PieceLength
	if pieceLength != 1048576 {
		t.Errorf("Expected info length to be 1048576 but was %v", pieceLength)
	}

	private := info.Private
	if private != 1 {
		t.Errorf("Expected private to be 1 but was %v", private)
	}

	directoryName := info.DirectoryName
	if directoryName != "" {
		t.Errorf("Expected no directory name for single file torrent but was %v", directoryName)
	}

	//Check files

	onlyFile := info.Files[0]
	if onlyFile == nil {
		t.Error("Failed to parse only file in single file torrent")
	}

	path := onlyFile.Path
	if path != "Reborn-OS-2018.09.17-x86_64.iso" {
		t.Errorf("Path is incorrect for single file %v", path)
	}

	fileLength := onlyFile.Length
	if fileLength != 1637744640 {
		t.Errorf("Expected only file length to be 1637744640 but should be %v", fileLength)
	}
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

	// Check basic metainfo

	announce := metaInfo.AnnounceUrls[0]
	if announce.String() != "http://legittorrents.info:2710/announce" {
		t.Errorf("Unexpected announce URL found for multi file torrent announce: %v", announce)
	}

	creationDate := metaInfo.CreationDate.Unix()
	if creationDate != 1536553238 {
		t.Errorf("Expected creation date to be 1536553238 but was %v", creationDate)
	}

	encoding := metaInfo.Encoding
	if encoding != "UTF-8" {
		t.Errorf("Expected encoding to be UTF-8 but was %v", encoding)
	}

	// Check info

	info := metaInfo.Info

	pieceLength := info.PieceLength
	if pieceLength != 524288 {
		t.Errorf("Expected info length to be 524288 but was %v", pieceLength)
	}

	private := info.Private
	if private != 0 {
		t.Errorf("Expected private to be 0 but was %v", private)
	}

	directoryName := info.DirectoryName
	if directoryName != "KJV" {
		t.Errorf("Expected directory name to be KJV for multi file torrent but was %v", directoryName)
	}

	//Check files [Checking first file and path for third only - there are a _lot_!]
	firstFile := info.Files[0]

	fileLength := firstFile.Length
	if fileLength != 9376812 {
		t.Errorf("Expected first file lenth to be 9376812 but was %v", fileLength)
	}

	filePath := firstFile.Path
	if filePath != "HolyBibleKJV.pdf" {
		t.Errorf("Expected first filepath to be HolyBibleKJV.pdf but was %v", filePath)
	}

	thirdFilePath := info.Files[2].Path
	if thirdFilePath != "Recordings/1 Chronicles 1.mp3" {
		t.Errorf("Expected third filepath to be 'Recordings/1 Chronicles 1.mp3' but was '%v'", thirdFilePath)
	}
}
