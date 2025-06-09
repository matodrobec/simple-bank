package gapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLoggerfunc(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	// log.Print("received a gRPC request")
	startTime := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}
	// log.Info().Str("protocol", "grpc").
	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		// Float64("duration_ms", float64(duration.Milliseconds())).
		Msg("received a gRPC request")
	return
}

type ResponseRecoreder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecoreder) Write(data []byte) (int, error) {
	r.Body = data
	return r.ResponseWriter.Write(data)
}

func (r *ResponseRecoreder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func HttpLogger(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		startTime := time.Now()
		// handler.ServeHTTP(res, req)
		recoreder := &ResponseRecoreder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(recoreder, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if recoreder.StatusCode >= 400 {
			err := fmt.Errorf(string(recoreder.Body))
			logger = log.Error().
				Err(err)
			// Bytes("body", recoreder.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", recoreder.StatusCode).
			Str("status_text", http.StatusText(recoreder.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
