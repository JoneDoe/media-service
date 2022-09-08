package models

import "os"

type LocalFile interface {
	Name() string
	OriginalName() string
	FileSystemPath() string
	MimeType() string
	SetPath(path string)
	Open() (*os.File, error)
	Remove() error
	ToJson() *MediaFile
	Uuid() string
}
