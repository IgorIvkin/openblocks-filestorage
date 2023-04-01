package main

import (
	api "openblocks/filestorage/api"
	app "openblocks/filestorage/application"

	"github.com/gin-gonic/gin"
)

func main() {
	application := app.NewApplication()
	defer application.CloseDb()

	router := gin.Default()
	// ограничиваем максимальный размер файла в 10 мегабайтов
	router.MaxMultipartMemory = 10 << 20

	router.POST("/api/v1/store", func(context *gin.Context) {
		api.ProcessStoreFile(application, context)
	})

	router.GET("/api/v1/file/:fileId", func(context *gin.Context) {
		api.ProcessGetFile(application, context)
	})

	router.Run(":8903")
}
