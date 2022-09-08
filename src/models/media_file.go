package models

import (
	"path/filepath"

	"github.com/mitchellh/mapstructure"
)

type MediaFile struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Version      string `json:"version"`
	OriginalName string `json:"originalName"`
	Uuid         string
}

func InitMedia(dbRecord interface{}) *MediaFile {
	media := &MediaFile{}
	mapstructure.Decode(dbRecord, media)

	return media
}

func (file *MediaFile) FileSystemPath() string {
	return filepath.Join(file.Path, file.Name)
}

func (file *MediaFile) Ext() string {
	return filepath.Ext(file.OriginalName)
}
