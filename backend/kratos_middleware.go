package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
)

type kratosMiddleware struct {
	ory *ory.APIClient
}

func NewMiddleware() *kratosMiddleware {
	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: os.Getenv("KRATOS_URL"), // Kratos Admin API
		},
	}
	return &kratosMiddleware{
		ory: ory.NewAPIClient(configuration),
	}
}

func (k *kratosMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("ory_kratos_session")

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "missing session cookie",
			})
			return
		}

		session, _, err := k.ory.FrontendAPI.
			ToSession(context.Background()).
			Cookie("ory_kratos_session=" + cookie).
			Execute()

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid session",
			})
			return
		}

		c.Set("user", session)

		c.Next()
	}
}
