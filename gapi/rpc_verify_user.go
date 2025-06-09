package gapi

import (
	"context"

	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/pb"
	"github.com/matodrobec/simplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEamilRequest) (*pb.VerifyEamilResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	verifyEmailResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	response := &pb.VerifyEamilResponse{
		IsVerified: verifyEmailResult.User.IsEmailVerified,
	}
	return response, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEamilRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidatePositiveNumber(req.EmailId); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	if err := validation.ValidateString(req.GetSecretCode(), 32, 128); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
