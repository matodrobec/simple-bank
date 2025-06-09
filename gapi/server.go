package gapi

import (
	"context"
	"errors"
	"fmt"

	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	db "github.com/matodrobec/simplebank/db/sqlc"
	_ "github.com/matodrobec/simplebank/doc/statik"
	"github.com/matodrobec/simplebank/pb"
	"github.com/matodrobec/simplebank/token"
	"github.com/matodrobec/simplebank/util"
	"github.com/matodrobec/simplebank/worker"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot craete token maker: %v", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}

func (server *Server) Start(
	ctx context.Context,
	waitGroup *errgroup.Group,
	address string,
) {

	grpcLoggerServerOption := grpc.UnaryInterceptor(GrpcLoggerfunc)
	grpcServer := grpc.NewServer(grpcLoggerServerOption)
	pb.RegisterSimpleBankServer(grpcServer, server)
	// for client documentation
	reflection.Register(grpcServer)

	// start server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC listener")
		// return fmt.Errorf("cannot start gRPC server at %s", address)
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
		err = grpcServer.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error().Err(err).Msg("gRPC server failed to server")
			return err
		}
		// if err != nil {
		// 	log.Error().Err(err).Msg("Cannot start grpc server")
		// 	return err
		// }
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})

	// log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	// return grpcServer.Serve(listener)
}

func (server *Server) StartProxy(
	ctx context.Context,
	waitGroup *errgroup.Group,
	address string) {
	// https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/#using-proto-names-in-json
	// use  snake_case naming convention
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	jsonPrettyOption := runtime.WithMarshalerOption("application/json+pretty", &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Indent:        "  ",
			Multiline:     true, // Optional, implied by presence of "Indent".
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// grpcMux := runtime.NewServeMux(jsonOption)
	grpcMux := runtime.NewServeMux(jsonOption, jsonPrettyOption)

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		// log.Fatal("canot register handler server")
		// return fmt.Errorf("canot register handler server")
		log.Fatal().Err(err).Msg("gRPC proxy: canot register handler server")
	}

	mux := http.NewServeMux()
	// all routes "/"
	mux.Handle("/", grpcMux)

	// doc/swagger
	// fileSwaggerServerHandler := http.FileServer(http.Dir("./doc/swagger"))
	fileSwaggerServerHandler, err := fs.NewWithNamespace("swagger")
	if err != nil {
		// return fmt.Errorf("cannot create statik file system: %s", err)
		log.Fatal().Err(err).Msg("cannot create static file system")
	}

	swaggerHnadler := http.StripPrefix("/swagger", http.FileServer(fileSwaggerServerHandler))
	mux.Handle("/swagger/", swaggerHnadler)

	// handlerWithCors := cors.Default().Handler(HttpLogger(mux))
	handlerWithCors := cors.
		New(cors.Options{
			// AllowedOrigins: []string{"http://localhost:8080", "http://test.com"},
			AllowedOrigins: server.config.CorsAllowedOrigins,
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
			},
			AllowCredentials: true,
		}).
		Handler(HttpLogger(mux))

	httpServer := &http.Server{
		// Handler: HttpLogger(mux),
		Handler: handlerWithCors,
		Addr:    address,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("HTTP gateway server failed to server")
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}
		log.Info().Msgf("HTTP gateway server is stopped")
		return nil
	})

	// // start server
	// listener, err := net.Listen("tcp", address)
	// if err != nil {
	// 	return fmt.Errorf("cannot start gRPC server at %s", address)
	// }

	// log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	// handlerLogger := HttpLogger(mux)
	// return http.Serve(listener, handlerLogger)
}
