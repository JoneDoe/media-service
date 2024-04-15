package local

import (
	"github.com/gin-gonic/gin"

	"istorage/attachment"
	"istorage/models"
	"istorage/services"
)

func Store(c *gin.Context, attach *models.Attachment) error {
	fm, err := attachment.FileManagerFactory(attachment.FileManagerConfig{
		MimeType: attach.OriginalFile.MimeType(),
		Version: "original",
	})

	if err != nil {
		return err
	}

	// Upload the file to specific dst.
	go func() {
		fm.SetFile(attach.OriginalFile)

		attach.Path = fm.DirManager().Path
		attach.Version = fm.ToJson().FileName

		c.SaveUploadedFile(attach.OriginalFile.Upload, fm.Filepath())
	}()

	return nil
}

func RemoveMedia(file *models.MediaFile, uuid string) error {
	if err := services.RemoveFile(file); err != nil {
		return err
	}

	return services.InitDb().DeleteRecord(uuid)
}

func ReadFile(file *models.MediaFile) string {
	return services.AbsolutePath(file)
}
