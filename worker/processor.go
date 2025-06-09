package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/mail"
	"github.com/matodrobec/simplebank/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	// Start() error
	Start(
		ctx context.Context,
		waitGroup *errgroup.Group,
	)
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmaiSender
	config util.Config
}

func NewRedisTaskProcessor(redisOpt asynq.RedisConnOpt, store db.Store, mailer mail.EmaiSender, config util.Config) TaskProcessor {

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					log.Error().Err(err).
						Str("type", task.Type()).
						Bytes("payload", task.Payload()).
						Msg("processed task faild")
				},
			),
			Logger: NewLogger(),

			// ErrorHandler: func() asynq.ErrorHandlerFunc {
			// 	return func(ctx context.Context, task *asynq.Task, err error) {
			// 	}
			// }(),
		},
	)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
		config: config,
	}
}

func (processor *RedisTaskProcessor) Start(
	ctx context.Context,
	waitGroup *errgroup.Group,
) {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	waitGroup.Go(func() error {
		log.Info().Msg("start task processor")
		err := processor.server.Start(mux)

		if err != nil {
			log.Error().Err(err).Msg("faild to start task processor")
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task prcessor")
		processor.server.Shutdown()

		log.Info().Msg("task processor is stopped")
		return nil
	})

}

// func (processor *RedisTaskProcessor) Start() error {
// 	mux := asynq.NewServeMux()
// 	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

// 	return processor.server.Start(mux)
// }
