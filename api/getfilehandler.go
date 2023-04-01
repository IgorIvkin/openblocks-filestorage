package api

import (
	"fmt"
	"log"
	"net/http"
	app "openblocks/filestorage/application"

	"github.com/gin-gonic/gin"
)

// Обрабатывает запрос на получение файла
func ProcessGetFile(application *app.Application, context *gin.Context) {
	fileId := context.Param("fileId")

	fileContent, fileType, err := application.GetFile(fileId)
	if err != nil {
		errorMessage := fmt.Sprintf("Cannot get file, reason: %v", err)
		log.Println(errorMessage)
		context.JSON(http.StatusNotFound, gin.H{
			"message": errorMessage,
		})
		return
	}

	context.Data(http.StatusOK, fileType, fileContent)
}
