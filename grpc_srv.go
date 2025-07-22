package srv

import (
	"context"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRPCConfig interface {
	Endpoint() string
	ApiEndpoint() string
}

type ProtoServiceServer interface {
	Register(grpcSrv *grpc.Server)
	RegisterFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}

type GRPCServer struct {
	logger *slog.Logger
	cfg    GRPCConfig
	ss     []ProtoServiceServer
}

func NewGRPCServer(logger *slog.Logger, cfg GRPCConfig, ss []ProtoServiceServer) *GRPCServer {
	return &GRPCServer{
		logger: logger,
		cfg:    cfg,
		ss:     ss,
	}
}

func (s *GRPCServer) Run() error {
	grpcListener, err := net.Listen("tcp", s.cfg.Endpoint())
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	if len(s.ss) > 0 {
		for _, svc := range s.ss {
			svc.Register(grpcServer)
		}
	}

	go func() {
		s.logger.Info("gRPC server running", slog.String("endpoint", s.cfg.Endpoint()))
		err := grpcServer.Serve(grpcListener)
		if err != nil {
			s.logger.Error("failed to serve gRPC", slog.Any("err", err))
		}
	}()
	return nil
}
