package models

import (
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type DownloadedFile struct {
	File       *os.File
	path, uuid string
}

func NewLocalFile(file *os.File) LocalFile {
	return &DownloadedFile{
		File: file,
		uuid: uuid.New().String(),
	}
}

func (d *DownloadedFile) SetPath(path string) {
	d.path = path
}

func (d *DownloadedFile) Name() string {
	return d.uuid + d.ext()
}

func (d *DownloadedFile) MimeType() string {
	return mime.TypeByExtension(d.ext())
}

func (d *DownloadedFile) FileSystemPath() string {
	return filepath.Join(d.path, d.Name())
}

func (d *DownloadedFile) Open() (*os.File, error) {
	return os.Open(d.File.Name())
}

func (d *DownloadedFile) Uuid() string {
	return d.uuid
}

func (d *DownloadedFile) Remove() error {
	return os.Remove(d.File.Name())
}

func (d *DownloadedFile) OriginalName() string {
	segments := strings.Split(d.File.Name(), "/")

	return segments[len(segments)-1]
}

func (d *DownloadedFile) ToJson() *MediaFile {
	return &MediaFile{
		Path:         d.path,
		Name:         d.Name(),
		OriginalName: d.OriginalName(),
		Type:         d.MimeType(),
		Uuid:         d.Uuid(),
		Version:      "1",
	}
}

func (d *DownloadedFile) ext() string {
	return filepath.Ext(d.File.Name())
}
