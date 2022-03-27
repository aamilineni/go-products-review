package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/aamilineni/go-products-review/constants"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// basic authentication middleware
		username, password, ok := c.Request.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(os.Getenv(constants.AUTH_USERNAME)))
			expectedPasswordHash := sha256.Sum256([]byte(os.Getenv(constants.AUTH_PASSWORD)))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				c.Next()
				return
			}
		}

		c.Writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		c.Abort()
	}
}
