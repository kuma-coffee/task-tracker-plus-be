package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		var t string

		t, err := ctx.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				contentType := ctx.Request.Header.Get("Content-type")
				if contentType == "" {
					contentType = "application/json"
				}

				authToken := ctx.Request.Header.Get("Authorization")
				if authToken == "" {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
					return
				}

				authTokenArr := strings.Fields(authToken)
				if len(authTokenArr) != 2 || authTokenArr[0] != "Bearer" {
					ctx.AbortWithStatus(http.StatusUnauthorized)
					return
				}

				t = authTokenArr[1]
			} else {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		claims := &model.Claims{}

		token, err := jwt.ParseWithClaims(t, claims, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		// ctx.Set("email", claims.Email)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Next()
	})
}
