package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Oabraham1/open-blogger/server/auth"
	"github.com/Oabraham1/open-blogger/server/log"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(authenticator auth.Authenticator) gin.HandlerFunc {
	genericError := gin.H{"error": "unauthorized"}
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if authorizationHeader == "" {
			logger.LogError(errors.New("authorization header is not provided"), "AuthenticationMiddleware")
			c.AbortWithStatusJSON(http.StatusUnauthorized, genericError)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			logger.LogError(errors.New("invalid authorization header format"), "AuthenticationMiddleware")
			c.AbortWithStatusJSON(http.StatusUnauthorized, genericError)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			logger.LogError(fmt.Errorf("unsupported authorization type %s", authorizationType), "AuthenticationMiddleware")
			c.AbortWithStatusJSON(http.StatusUnauthorized, genericError)
			return
		}

		token := fields[1]
		payload, err := authenticator.VerifyToken(token)
		if err != nil {
			logger.LogError(err.Error(), "AuthenticationMiddleware")
			c.AbortWithStatusJSON(http.StatusUnauthorized, genericError)
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}
