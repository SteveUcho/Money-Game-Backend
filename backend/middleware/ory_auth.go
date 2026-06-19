package middleware

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
	"steveucho.com/packages/backend/models"
)

type kratosMiddleware struct {
	ory *ory.APIClient
}

func NewAuthMiddleware() *kratosMiddleware {
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

func GetTraits(session *ory.Session) (*models.IdentityTraits, error) {
	var traits models.IdentityTraits
	b, err := json.Marshal(session.Identity.Traits)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &traits)
	if err != nil {
		return nil, err
	}
	return &traits, nil
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

		traits, err := GetTraits(session)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error": "failed to get user traits",
			})
			return
		}

		c.Set("user", &models.User{
			Session: session,
			Traits:  *traits,
		})
		c.Next()
	}
}
