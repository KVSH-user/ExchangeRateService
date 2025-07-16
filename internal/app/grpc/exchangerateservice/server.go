package exchangerateservice

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	pb "github.com/KVSH-user/ExchangeRateService/pkg/pb/exchangerateservice"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	logger     *slog.Logger
	port       string

	exchangeRateModule *ExchangeRateService
}

func NewServer(log *slog.Logger, port string, exchangeRateModule ExchangeRateModule) *Server {
	if log == nil {
		log = slog.Default()
	}

	return &Server{
		logger:             log,
		port:               port,
		exchangeRateModule: NewExchangeRateService(log, exchangeRateModule),
	}
}

func (s *Server) Start(_ context.Context) error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	s.listener = listener

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.loggingInterceptor(),
			s.recoveryInterceptor(),
		),
	)

	s.registerServices()

	reflection.Register(s.grpcServer)

	s.logger.Info("Starting gRPC server", "port", s.port)

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop(ctx context.Context) error {
	if s.grpcServer == nil {
		return nil
	}

	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("Force stopping gRPC server")
		s.grpcServer.Stop()
		return ctx.Err()
	}
}

func (s *Server) registerServices() {
	pb.RegisterExchangeRateServiceServer(s.grpcServer, s.exchangeRateModule)
}

func (s *Server) loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		logLevel := slog.LevelInfo
		if err != nil {
			logLevel = slog.LevelError
		}

		s.logger.Log(ctx, logLevel, "gRPC call",
			"method", info.FullMethod,
			"duration", duration,
			"error", err,
		)

		return resp, err
	}
}

func (s *Server) recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("gRPC panic recovered",
					"method", info.FullMethod,
					"panic", r,
				)
				err = fmt.Errorf("error code %d - %s", codes.Internal, r.(string))
			}
		}()

		return handler(ctx, req)
	}
}
