package api

import (
	"fmt"
	"log"
	"net/http"
	app "openblocks/filestorage/application"

	"github.com/gin-gonic/gin"
)

// Обрабатывает запрос на получение файла
func ProcessGetFile(application *app.Application, ctx *gin.Context) {
	fileId := ctx.Param("fileId")

	fileContent, fileType, err := application.GetFile(fileId)
	if err != nil {
		errorMessage := fmt.Sprintf("Cannot get file, reason: %v", err)
		log.Println(errorMessage)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": errorMessage,
		})
		return
	}

	ctx.Data(http.StatusOK, fileType, fileContent)
}
