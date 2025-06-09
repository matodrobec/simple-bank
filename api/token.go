package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/matodrobec/simplebank/db/sqlc"
)

type renewTokenRequest struct {
	RefresToken string `json:"refres_token" binding:"required"`
}

type renewTokenResponse struct {
	// SessionID            uuid.UUID          `json:"session_id"`
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken         string    `json:"refres_token"`
	AccessRefreshToken   time.Time `json:"refres_token_expires_at"`
	// User                 createUserResponse `json:"user"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var req renewTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify token
	refresPayload, err := server.tokenMaker.VerifyToken(req.RefresToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// session from db
	session, err := server.store.GetSession(ctx, refresPayload.ID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, db.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refresPayload.Username {
		err := fmt.Errorf("incored session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if req.RefresToken != session.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, refresPayload.Username)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, db.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, errorResponse(err))
		return
	}

	// access token
	accessToken, accessPayload, err := server.tokenMaker.CrateToken(
		refresPayload.Username,
		server.config.AccessTokenDuration,
		user.Role,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// refres token
	refreshToken, refresPayload, err := server.tokenMaker.CrateToken(
		refresPayload.Username,
		server.config.RefresTokenDuration,
		user.Role,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sessionArgs := db.RotateSessionTxParams{
		FromSessionID: session.ID,
		ToSession: db.CreateSessionParams{
			ID:           refresPayload.ID,
			Username:     refresPayload.Username,
			RefreshToken: refreshToken,
			UserAgent:    ctx.Request.UserAgent(),
			ClientIp:     ctx.ClientIP(),
			IsBlocked:    false,
			ExpiresAt:    refresPayload.ExpiredAt,
		},
	}
	if _, err := server.store.RotateSessionTx(ctx, sessionArgs); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewTokenResponse{
		// SessionID:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		RefreshToken:         refreshToken,
		AccessRefreshToken:   refresPayload.ExpiredAt,
		// User:                 newUserRespnse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
