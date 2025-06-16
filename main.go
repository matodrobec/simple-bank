package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	// _ "github.com/jackc/pgx/v5"
	// _ "github.com/lib/pq"
	"github.com/matodrobec/simplebank/api"
	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/gapi"
	"github.com/matodrobec/simplebank/mail"
	"github.com/matodrobec/simplebank/util"
	"github.com/matodrobec/simplebank/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServerType int

const (
	ServerGin ServerType = iota
	ServerGRPC
	ServerGRPCProxy
)

var typeName = map[ServerType]string{
	ServerGin:       "gin",
	ServerGRPC:      "grpc",
	ServerGRPCProxy: "proxy",
}

func (st ServerType) String() string {
	return typeName[st]
}

func ParseServerType(s string) (st ServerType) {
	st = ServerGin
	for v, k := range typeName {
		if s == k {
			st = v
			return
		}
	}
	return
}

var interruptSignals = []os.Signal{
	os.Interrupt,
	// Sent when the user presses Ctrl+C in the terminal.
	syscall.SIGTERM,
	// A termination request from the OS or another process (e.g., Kubernetes shutdown).
	syscall.SIGINT,
}

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://postgres:test@localhost:5432/bank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {

	// err := util.LoadConfigAndWatcing(".", func(c util.Config) {
	// 	runServer(c)
	// })

	// if err != nil {
	// 	log.Fatal().Msg("cannot load config:", err)
	// }

	config, err := util.LoadConfig(".")

	if config.Environment == util.DevEnv {
		// Pretty logging to console
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	// connPool, err := pgxpool.New(context.Background(), config.DBSource)
	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot conntect to db")
	}

	runDBMigrate(config.MigrationUrl, config.DBSource)
	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	// runTaskProcessor(ctx, waitGroup, redisOpt, store, config)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}

	// go runTaskProcessor(redisOpt, store, config)
	// go runGatewayServer(config, store, taskDistributor)
	// runGrpcServer(config, store, taskDistributor)

	// var stArg string
	// if len(os.Args) > 1 {
	// 	stArg = os.Args[1]
	// }
	// st := ParseServerType(stArg)

	// server, err := NewServerFactory(st, config, store)
	// if err != nil {
	// 	log.Fatal().Msg("cannot create server:", err)
	// }
	// err = server.Start(config.ServerAddress)
	// if err != nil {
	// 	log.Fatal().Msg("Cannot start server: ", err)
	// }

}

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	config util.Config,
) {
	// mailer := mail.NewGenericSender(config)
	// mailer := mail.NewSmtpSender(config)
	mailer := mail.NewGmailSender(config.GetFromName(), config.GetFromEmailAddress(), config.GetSmtpPassword())
	runTaskProcessor := worker.NewRedisTaskProcessor(
		redisOpt, store, mailer, config,
	)

	runTaskProcessor.Start(ctx, waitGroup)

	// err := runTaskProcessor.Start()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("faild to start task processor")
	// } else {
	// 	log.Info().Msg("start task processor")
	// }
}

func runDBMigrate(migrationUrl, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal().Msgf("cannot create new migrate instance: %s", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("failed to run migrate up: %s", err)
	}
	log.Print("db migreate sucessfully")
}

// type RunServer interface {
// 	Start(address string) error
// }

// func NewServerFactory(serverType ServerType, config util.Config, store db.Store) (server RunServer, err error) {
// 	switch serverType {
// 	case ServerGRPC:
// 		server, err = gapi.NewServer(config, store)
// 	default:
// 		server, err = api.NewServer(config, store)
// 	}

// 	return
// }

func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config, store db.Store,
	taskDistributor worker.TaskDistributor,
) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create grpc server")
	}

	server.StartProxy(ctx, waitGroup, config.HTTPServerAddress)

	// err = server.StartProxy(config.HTTPServerAddress)
	// if err != nil {
	// 	log.Fatal().Msg("Cannot start HTTP gateway server")
	// }
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create grpc server")
	}

	server.Start(ctx, waitGroup, config.GRPCServerAddress)

	// err = server.Start(config.GRPCServerAddress)
	// if err != nil {
	// 	log.Fatal().Msg("Cannot start grpc server")
	// }

	// server, err := gapi.NewServer(config, store)
	// if err != nil {
	// 	log.Fatal().Msg("cannot create server:", err)
	// }

	// grpcServer := grpc.NewServer()
	// pb.RegisterSimpleBankServer(grpcServer, server)
	// // for client documentation
	// reflection.Register(grpcServer)

	// // start server
	// listener, err := net.Listen("tcp", config.GRPCServerAddress)
	// if err != nil {
	// 	log.Fatal().Msg("cannot start gRPC server at ", config.GRPCServerAddress)
	// }

	// log.Printf("start gRPC server at %s", listener.Addr().String())
	// err = grpcServer.Serve(listener)
	// if err != nil {
	// 	log.Fatal().Msg("cannot start gRPC server")
	// }
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create gin server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("Cannot start gin server")
	}
}

// func runServer(config util.Config) {
// 	conn, err := sql.Open(config.DBDriver, config.DBSource)
// 	if err != nil {
// 		log.Fatal().Msg("cannot conntect to db: ", err)
// 	}

// 	store := db.NewStore(conn)

// 	server := api.NewServer(store)
// 	err = server.Start(config.ServerAddress)
// 	if err != nil {
// 		log.Fatal().Msg("Cannot start server: ", err)
// 	}
// }
