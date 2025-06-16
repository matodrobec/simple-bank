package gapi

import (
	"context"

	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/pb"
	"github.com/matodrobec/simplebank/util"
	"github.com/matodrobec/simplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	// log.Println(">> create user start")
	// time.Sleep(20 * time.Second)

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	txArg := db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.User) error {
			return nil
		},
		// AfterCreate: func(user db.User) error {
		// 	opts := []asynq.Option{
		// 		asynq.MaxRetry(10),
		// 		// we ned this time because transaction is commit after AfterCraete
		// 		// there is transaction commit
		// 		// without 10s delay can happen that user is not commited into DB
		// 		asynq.ProcessIn(2 * time.Second),
		// 		asynq.Queue(worker.QueueCritical),
		// 	}
		// 	taskPayload := &worker.PayloadSendVerifyEmail{
		// 		Username: user.Username,
		// 	}
		// 	return server.taskDistributor.DistributedTaskSendEmail(ctx, taskPayload, opts...)
		// },
	}

	// user, err := server.store.CreateUser(ctx, arg)
	txResult, err := server.store.CreateUserTx(ctx, txArg)
	if err != nil {
		// if pqErr, ok := err.(*pq.Error); ok {
		// 	switch pqErr.Code.Name() {
		// 	case "unique_violation":
		// 		return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
		// 	}
		// }
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err)
			// return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "faild to craete user: %s", err)
	}

	// TODO: use db transaction
	// opts := []asynq.Option{
	// 	asynq.MaxRetry(10),
	// 	asynq.ProcessIn(10 * time.Second),
	// 	asynq.Queue(worker.QueueCritical),
	// }
	// taskPayload := &worker.PayloadSendVerifyEmail{
	// 	Username: user.Username,
	// }
	// err = server.taskDistributor.DistributedTaskSendEmail(ctx, taskPayload, opts...)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to distribute task to send verify email: %s", err)
	// }
	response := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}

	// log.Println(">> create user done")

	return response, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := validation.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err)) // full_name see. rpc_create_user.proto CreateUserRequest
	}
	if err := validation.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err)) // email see. rpc_create_user.proto
	}
	return violations
}
