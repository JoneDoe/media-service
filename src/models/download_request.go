package models

import (
	"encoding/json"
	"net/url"
	"strings"
)

type DownloadRequestItem struct {
	Source           string `json:"source" binding:"required,url"`
	ThumbnailProfile string `json:"thumbnail-profile"`
}

func (r *DownloadRequestItem) BuildDestinationFileName() string {
	fileUrl, _ := url.Parse(r.Source)

	segments := strings.Split(fileUrl.Path, "/")

	return segments[len(segments)-1]
}

type DownloadRequest struct {
	Items []DownloadRequestItem `binding:"required,dive"`
}

func (i *DownloadRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &i.Items)
}
