package gapi

import (
	"context"
	"errors"

	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/pb"
	"github.com/matodrobec/simplebank/util"
	"github.com/matodrobec/simplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		// status := http.StatusInternalServerError

		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "faild to find user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorect password")

	}

	accessToken, accessPayload, err := server.tokenMaker.CrateToken(
		req.Username,
		server.config.AccessTokenDuration,
		user.Role,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "faild to create access token")
	}

	refreshToken, refresPayload, err := server.tokenMaker.CrateToken(
		req.Username,
		server.config.RefresTokenDuration,
		user.Role,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "faild to create refres token")
	}

	mtd := extractMetadata(ctx)

	sessionArgs := db.CreateSessionParams{
		ID:           refresPayload.ID,
		Username:     refresPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtd.UseAgnet,
		ClientIp:     mtd.ClientIp,
		// UserAgent:    ctx.Request.UserAgent(),
		// ClientIp:     ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refresPayload.ExpiredAt,
	}
	session, err := server.store.CreateSession(ctx, sessionArgs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "faild to create session")
	}

	rsp := &pb.LoginUserResponse{
		SessionId:            session.ID.String(),
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:         refreshToken,
		RefresTokenExpiresAt: timestamppb.New(refresPayload.ExpiredAt),
		User:                 convertUser(user),
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err, "LOGIN_PASSWORD"))
	}
	return violations
}
