package imaginary

import (
	"errors"
	"os"
	"strings"

	"istorage/models"
)

type ResizeProfile struct {
	ProfileName string `uri:"profile" binding:"required"`
	MediaFile   *models.MediaFile
}

type ResizeResult struct {
	FilePath, FileName string
}

func (p *ResizeProfile) Resize(filePath string) (*ResizeResult, error) {
	settings, err := loadProfile(p.ProfileName)
	if err != nil {
		return nil, errors.New(strings.Join([]string{
			"Can`t make operation, try one of following: ",
			AvailableProfiles(),
		}, ""))
	}

	pattern := strings.Join([]string{"cropper", ".*", p.MediaFile.Ext()}, "")

	tmpFile, _ := os.CreateTemp("", pattern)

	Resize(settings, filePath, tmpFile.Name())

	return &ResizeResult{
		FileName: "resized" + p.MediaFile.Ext(),
		FilePath: tmpFile.Name(),
	}, nil
}
