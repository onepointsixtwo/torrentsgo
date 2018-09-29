package model

import (
	"net/url"
	"time"
)

// TYPES

type MetaInfo struct {
	AnnounceUrls []url.URL
	CreationDate time.Time
	Comment      string
	CreatedBy    string
	Encoding     string
	Info         *Info
}

/*
	Note: I've collapsed down the spec to basically just do one type of info struct for single or multiple files.
	This way the directory name can be blank and the path can just be the filename for single file, and it
	keeps it as only one 'type' that has to be handled which just seems simpler.
*/

type Info struct {
	PieceLength   int
	Pieces        []byte
	Private       int
	Files         []*File
	DirectoryName string
}

type File struct {
	Path   string
	Length int
	Md5Sum string
}

// INITIALISATION

func NewMetaInfo(announceUrls []url.URL,
	creationDate time.Time,
	comment string,
	createdBy string,
	encoding string,
	info *Info) *MetaInfo {
	return &MetaInfo{announceUrls, creationDate, comment, createdBy, encoding, info}
}

func NewInfo(pieceLength int,
	pieces []byte,
	private int,
	files []*File,
	directoryName string) *Info {
	return &Info{pieceLength, pieces, private, files, directoryName}
}

func NewFile(path string, length int, md5Sum string) *File {
	return &File{path, length, md5Sum}
}
