package models

type UploadContext struct {
	Context string `uri:"context" binding:"required"`
}
