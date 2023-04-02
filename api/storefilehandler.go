package api

import (
	"fmt"
	"io"
	"log"
	"net/http"

	app "openblocks/filestorage/application"

	"github.com/gin-gonic/gin"
)

// Обрабатывает запрос на сохранение нового файла
func ProcessStoreFile(application *app.Application, ctx *gin.Context) {

	fileType := ctx.Request.FormValue("file_type")
	if fileType == "" {
		log.Println("Cannot store file, file type is missing")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot store file, file type is missing",
		})
		return
	}

	formFile, _ := ctx.FormFile("file_to_store")
	if formFile == nil {
		log.Println("Cannot store file, nothing is uploaded")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot store file, nothing is uploaded",
		})
		return
	}

	file, err := formFile.Open()
	if err != nil {
		log.Printf("Cannot process file %s, reason: %v\n", formFile.Filename, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot process uploaded file",
		})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Cannot read content of uploaded file %s, reason: %v\n", formFile.Filename, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot read content of uploaded file",
		})
		return
	}

	volume := application.ChooseIdleVolume()
	fileId, err := volume.StoreFile(content, fileType)
	if err != nil {
		message := fmt.Sprintf("Cannot store file %s, reason: %v", formFile.Filename, err)
		log.Println(message)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"file": fileId,
		})
	}
}
