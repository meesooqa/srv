package srv

import (
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type GRPCConfig interface {
	Endpoint() string
}

type ProtoServiceServer interface {
	Register(grpcSrv *grpc.Server)
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
