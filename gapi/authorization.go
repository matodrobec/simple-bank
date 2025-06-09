package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/matodrobec/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	grpcGatewayAuthorizationHeader = "grpcgateway-authorization"
	authorizationHader             = "authorization"
	authorizationBearer            = "bearer"
)

func (server *Server) autohorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHader)
	if len(values) == 0 {
		values = md.Get(grpcGatewayAuthorizationHeader)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	authHedersFields := strings.Fields(authHeader)
	if len(authHedersFields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(authHedersFields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsaported authorization type %s", authType)
	}

	accessToken := authHedersFields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	if !hasPermissions(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermissions(userRole string, accesibleRolles []string) bool {
	for _, role := range accesibleRolles {
		if role == userRole {
			return true
		}
	}
	return true
}
