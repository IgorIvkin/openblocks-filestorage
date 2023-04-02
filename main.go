package main

import (
	api "openblocks/filestorage/api"
	app "openblocks/filestorage/application"
	middleware "openblocks/filestorage/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	application := app.NewApplication()
	defer application.CloseDb()

	router := gin.Default()

	// Задаём middleware для проверки JWT-токена в приложении
	router.Use(middleware.AuthenticateWithJwt(application))

	// Ограничиваем максимальный размер файла в 10 мегабайтов
	router.MaxMultipartMemory = 10 << 20

	router.POST("/api/v1/store", func(ctx *gin.Context) {
		api.ProcessStoreFile(application, ctx)
	})

	router.GET("/api/v1/file/:fileId", func(ctx *gin.Context) {
		api.ProcessGetFile(application, ctx)
	})

	router.Run(":8903")
}
