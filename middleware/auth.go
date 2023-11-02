package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		c, err := ctx.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				contentType := ctx.Request.Header.Get("Content-type")
				if contentType == "" {
					ctx.AbortWithStatus(http.StatusSeeOther)
					return
				}
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		claims := &model.Claims{}

		token, err := jwt.ParseWithClaims(c, claims, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		ctx.Set("email", claims.Email)

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
