package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"istorage/models"
	"istorage/services"
	"istorage/utils"
)

type FileProwlerController utils.DatabaseStructure

func (f *FileProwlerController) ProwlFile(c *gin.Context) {
	var request models.DownloadRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		utils.Response{c}.Error(http.StatusUnprocessableEntity, "Invalid request format")

		return
	}

	c.JSON(http.StatusCreated, services.Prowler(f.DB, request.Items))
}
