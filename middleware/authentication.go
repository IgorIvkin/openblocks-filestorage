package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	app "openblocks/filestorage/application"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// В том случае, если в приложении требуется проверять JWT-токены, проверяет JWT-токен по
// указанным ключам JWKS, если токен требуется проверить, и он не валидный, немедленно
// происходит прекращение обработки запроса с ошибкой 401.
func AuthenticateWithJwt(application *app.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var authenticationFailed bool = false

		// Проверяем JWT-токен:
		// 1. если в настройках приложения задано, что нужно использовать JWT-аутентификацию
		// 2. если текущий URL запроса не находится в списке URL на исключение верификации JWT - `general.except-urls`
		if application.Config.General.UseJwtAuth && !isUrlExceptedToAuthenticate(application, ctx) {
			strToken := getBearerToken(ctx)
			jwks := fetchJwks(application.Config.General.JwksUrl)
			if strToken != "" {

				// Парсим токен, используя JWKS-ключи, которые мы получили из Keycloak,
				// если токен распарсить не получается, возвращается ошибка 401
				jwtToken, err := jwt.ParseString(strToken, jwt.WithKeySet(jwks))
				if err != nil {
					authenticationFailed = true
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"message": fmt.Sprintf("Cannot parse JWT-token, reason: %v", err),
					})
				}

				// Верифицируем токен, если токен некорректный, возвращается ошибка 401
				if err == nil {
					err := jwt.Validate(jwtToken)
					if err != nil {
						authenticationFailed = true
						ctx.JSON(http.StatusUnauthorized, gin.H{
							"message": fmt.Sprintf("JWT-token is invalid, reason: %v", err),
						})
					}
				}

			} else {

				// Если вообще не задан JWT-токен, а использовать его требуется по конфигурации
				// приложения, возвращается ошибка 401
				authenticationFailed = true
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": "No JWT-token provided, request declined",
				})
			}
		}

		// Продолжаем выполнение запроса только в том случае, если аутентификация
		// не требуется или не провалена
		if authenticationFailed {
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

func isUrlExceptedToAuthenticate(application *app.Application, ctx *gin.Context) bool {
	urlPath := ctx.Request.URL.Path
	exceptUrls := application.Config.General.ExceptUrls
	for _, exceptUrl := range exceptUrls {
		if strings.HasPrefix(urlPath, exceptUrl) {
			return true
		}
	}
	return false
}

func fetchJwks(jwksUrl string) jwk.Set {
	jwks, err := jwk.Fetch(context.Background(), jwksUrl)
	if err != nil {
		log.Fatalf("Cannot fetch JWKS to validate JWT-token, reason: %v", err)
	}
	return jwks
}

func getBearerToken(ctx *gin.Context) string {
	const BEARER = "Bearer "
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader != "" {
		return authorizationHeader[len(BEARER):]
	} else {
		return ""
	}
}
