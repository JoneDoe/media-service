package models

type DownloadResponse OutputModel

type DownloadResponseError struct {
	Source  string `json:"source"`
	Message string `json:"message"`
}
