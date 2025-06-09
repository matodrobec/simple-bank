package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/token"
	"github.com/matodrobec/simplebank/util"
	"github.com/matodrobec/simplebank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymetricKey:    util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}

func newContextWithBeareToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) context.Context {

	token, payload, err := tokenMaker.CrateToken(username, duration, role)
	require.NoError(t, err)
	require.NotNil(t, payload)

	bareToken := fmt.Sprintf(
		"%s %s",
		authorizationBearer,
		token,
	)

	md := metadata.MD{
		authorizationHader: []string{
			bareToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
