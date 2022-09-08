package models

import (
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

const FileTypeImage = "image"

// Attachment contain info about original uploaded file, uuid...
type Attachment struct {
	OriginalFile              *OriginalFile
	Path, Version, Uuid, Name string
}

// Return Attachment
func Create(file *multipart.FileHeader) *Attachment {
	uuid := uuid.New().String()
	originalFile := &OriginalFile{file}

	return &Attachment{
		OriginalFile: originalFile,
		Uuid:         uuid,
		Name:         uuid + originalFile.Ext(),
	}
}

func (attachment *Attachment) ToJson() *MediaFile {
	return &MediaFile{
		Path:         attachment.Path,
		Name:         attachment.Name,
		OriginalName: attachment.OriginalFile.Filename(),
		Type:         attachment.OriginalFile.MimeType(),
		Version:      attachment.Version,
		Uuid:         attachment.Uuid,
	}
}

func (attachment *Attachment) FileSystemPath() string {
	return filepath.Join(attachment.Path, attachment.Name)
}
