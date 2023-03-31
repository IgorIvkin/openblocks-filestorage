package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	application := NewApplication()
	defer application.CloseDb()

	router := gin.Default()
	// ограничиваем максимальный размер файла в 10 мегабайтов
	router.MaxMultipartMemory = 10 << 20

	router.POST("/api/v1/store", func(context *gin.Context) {
		ProcessStoreFile(application, context)
	})

	router.Run(":8903")
}
