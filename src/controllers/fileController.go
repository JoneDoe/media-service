package controllers

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"istorage/imaginary"
	"istorage/logger"
	"istorage/models"
	"istorage/repositories/cloud"
	"istorage/utils"
)

type FileController utils.DatabaseStructure

var (
	dummyPath  = "default/00/00/00"
	dummyImage = "./data/dummy.png"
)

func (f *FileController) ReadFile(c *gin.Context) {
	token, _ := bindRequestToken(c)

	resp := utils.Response{Context: c}

	media, err := f.findMediaByToken(token)
	if err != nil {
		resp.Error(http.StatusNotFound, err.Error())

		return
	}

	if media.Path == dummyPath {
		c.File(dummyImage)

		return
	}

	/*if err = services.Check(media); err != nil {
		resp.Error(http.StatusNotFound, "Can`t read file")

		return
	}*/

	resizeProfile, _ := bindResizeProfile(c)
	if resizeProfile.ProfileName != "" {
		resizeProfile.MediaFile = media
		c.Set("resizeProfile", resizeProfile)

		return
	}

	filePath := cloud.ReadFile(media)
	defer os.Remove(filePath)

	if media.Type == models.FileTypeImage {
		c.File(filePath)
	} else {
		c.FileAttachment(filePath, media.OriginalName)
	}
}

func (f *FileController) ReadFileWithResize(c *gin.Context) {
	resizeProfile := c.MustGet("resizeProfile").(*imaginary.ResizeProfile)

	resp := utils.Response{Context: c}

	if resizeProfile.MediaFile.Type != models.FileTypeImage {
		resp.Error(http.StatusUnsupportedMediaType, "Operation allowed only for image files")

		return
	}

	filePath := cloud.ReadFile(resizeProfile.MediaFile)
	defer os.Remove(filePath) // clean up

	out, err := resizeProfile.Resize(filePath)
	defer os.Remove(out.FilePath) // clean up

	if err != nil {
		resp.Error(http.StatusNotFound, err.Error())

		return
	}

	c.FileAttachment(out.FilePath, out.FileName)
}

func (f *FileController) DeleteFile(c *gin.Context) {
	token, _ := bindRequestToken(c)

	resp := utils.Response{Context: c}

	media, err := f.findMediaByToken(token)
	if err != nil {
		resp.Error(http.StatusNotFound, err.Error())

		return
	}

	go func() {
		if err = cloud.RemoveMedia(media); err != nil {
			logger.Error(err)
		}

		f.DB.DeleteRecord(token.Uuid)
	}()

	resp.Success(http.StatusOK, token.Uuid)
}

func (f *FileController) FileInfo(c *gin.Context) {
	token, _ := bindRequestToken(c)

	media, err := f.findMediaByToken(token)
	if err != nil {
		utils.Response{Context: c}.Error(http.StatusNotFound, err.Error())

		return
	}

	utils.Response{Context: c}.Success(http.StatusOK, models.MetaData{
		OutputModel: models.OutputModel{
			FileName: media.OriginalName,
			Uuid:     token.Uuid,
		},
		Key: media.Name,
		Url: cloud.GetFullUrl(media.FileSystemPath()),
	})
}

func (f *FileController) findMediaByToken(token *models.RequestToken) (*models.MediaFile, error) {
	rec, err := f.DB.GetRecord(token.Uuid)
	if err != nil {
		return nil, errors.New("record not found by token " + token.Uuid)
	}

	return models.InitMedia(rec), nil
}

func bindRequestToken(c *gin.Context) (*models.RequestToken, error) {
	token := &models.RequestToken{}
	if err := c.ShouldBindUri(&token); err != nil {
		return nil, err
	}

	return token, nil
}

func bindResizeProfile(c *gin.Context) (*imaginary.ResizeProfile, error) {
	data := &imaginary.ResizeProfile{}
	if err := c.ShouldBindUri(&data); err != nil {
		return nil, err
	}

	return data, nil
}
