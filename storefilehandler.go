package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Обрабатывает запрос на сохранение нового файла
func ProcessStoreFile(application *Application, context *gin.Context) {

	formFile, _ := context.FormFile("file_to_store")

	file, err := formFile.Open()
	if err != nil {
		log.Printf("Cannot process file %s, reason: %v\n", formFile.Filename, err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot process uploaded file",
		})
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Cannot read content of uploaded file %s, reason: %v\n", formFile.Filename, err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot read content of uploaded file",
		})
	}

	volume := application.ChooseIdleVolume()
	fileId, err := volume.StoreFile(content)
	if err != nil {
		message := fmt.Sprintf("Cannot store file %s, reason: %v", formFile.Filename, err)
		log.Println(message)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"file": fileId,
		})
	}
}
