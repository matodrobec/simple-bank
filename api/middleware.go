package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matodrobec/simplebank/token"
)

const (
    authorizationHeaderKey = "authorization"
    authorizationTypeBearer = "bearer"
    authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
        authrizationHeader := ctx.GetHeader(authorizationHeaderKey)
        if len(authrizationHeader) == 0 {
            err := errors.New("authorization header is not provided")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
            return
        }

        authorizationFields := strings.Fields(authrizationHeader)
        if len(authorizationFields) < 2 {
            err := errors.New("invalid authorization header format")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
            return
        }

        authorizationType := strings.ToLower(authorizationFields[0])
        if authorizationType != authorizationTypeBearer {
            err := fmt.Errorf("unsupported authorization type %s", authorizationType)
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
            return
        }

        accessToken := authorizationFields[1]
        payload, err := tokenMaker.VerifyToken(accessToken)
        if err != nil {
            // err := errors.New("invalid access token")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
            return
        }
        ctx.Set(authorizationPayloadKey, payload)
        ctx.Next()
	}
}
