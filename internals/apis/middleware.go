package apis

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

const (
	ORG_CONTEXT_NAME = "x-organization-id"
)

type TerraxenClaims struct {
	jwt.StandardClaims
	OrganizationId string
	UserId         string
}

type Middlewares struct {
	apiSvc services.ApiKeyProvider
}

func NewMiddlewares(apiKeySvc services.ApiKeyProvider) *Middlewares {
	mdw := &Middlewares{
		apiSvc: apiKeySvc,
	}
	return mdw
}

func (m *Middlewares) ValidateRequest(c *gin.Context) {
	apiKey := c.GetHeader("X-Terraxen-API")
	authHeader := c.GetHeader("authorization")
	if apiKey != "" {
		ak := m.apiSvc.GetKey(apiKey)
		if ak == nil {
			responseError(c, http.StatusForbidden, "Api key rejected")
			return
		}
		c.Set(ORG_CONTEXT_NAME, ak.OrganizationId)
		c.Next()
		return
	}

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")

		m.validateJWT(parts[1], []byte("terraxen"))
	}

	c.Set(ORG_CONTEXT_NAME, "")
	c.Next()
}

func (m *Middlewares) validateJWT(jwtToken string, key []byte) {

	token, err := jwt.ParseWithClaims(jwtToken, &TerraxenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return key, nil
	})
	if err != nil {
		log.Fatalf("JWT verify failed: %v", err)
	}

	claims := token.Claims.(*TerraxenClaims)
	fmt.Println(claims.UserId)
}
