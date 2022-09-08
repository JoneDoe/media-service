package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"istorage/dir"
	"istorage/logger"
	"istorage/models"
	"istorage/repositories/cloud"
	"istorage/utils"
)

type AttachmentController utils.DatabaseStructure

func (a *AttachmentController) StoreAttachment(c *gin.Context) {
	files, err := validateRequest(c)

	if err != nil {
		utils.Response{c}.Error(http.StatusBadRequest, err.Error())

		return
	}

	filesList := a.processUpload(files, "")

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "files": filesList})
}

func (a *AttachmentController) StoreAttachmentWithContext(c *gin.Context) {
	resp := utils.Response{c}

	data := models.UploadContext{}
	if err := c.ShouldBindUri(&data); err != nil {
		resp.Error(http.StatusNotFound, err.Error())

		return
	}

	files, err := validateRequest(c)

	if err != nil {
		resp.Error(http.StatusBadRequest, err.Error())

		return
	}

	filesList := a.processUpload(files, data.Context)

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "files": filesList})
}

func (a *AttachmentController) processUpload(files []*multipart.FileHeader, context string) []models.OutputModel {
	var filesList = make([]models.OutputModel, 0)

	if context == "" {
		context = "default"
	}

	for _, file := range files {
		attach := models.Create(file)

		fm, _ := dir.NewManager(&dir.Config{
			MimeType: attach.OriginalFile.MimeType(),
			FileName: attach.Name,
			Context:  context,
		})

		attach.Path = fm.Abs()

		go func() {
			if err := cloud.Store(attach); err != nil {
				logger.Error(err)
			}

			if err := a.DB.CreateRecord(attach.ToJson()); err != nil {
				logger.Error(err)
			}
		}()

		filesList = bindAttachedResponse(filesList, attach)
	}

	return filesList
}

func bindAttachedResponse(list []models.OutputModel, attach *models.Attachment) []models.OutputModel {
	return append(list, models.OutputModel{
		FileName: attach.OriginalFile.Filename(),
		Uuid:     attach.Uuid,
		Url:      cloud.GetFullUrl(attach.FileSystemPath()),
	})
}

func validateRequest(c *gin.Context) ([]*multipart.FileHeader, error) {
	errMsg := "Undefined any file"

	if c.GetHeader("Content-Length") == "0" {
		return nil, errors.New(errMsg)
	}

	form, _ := c.MultipartForm()

	files := form.File["files[]"]
	if len(files) == 0 {
		return nil, errors.New(errMsg)
	}

	return files, nil
}
