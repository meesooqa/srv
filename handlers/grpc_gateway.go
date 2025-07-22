package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/meesooqa/srv"
)

type GrpcGateway struct {
	logger *slog.Logger
	cfg    srv.GRPCConfig
	ss     []srv.ProtoServiceServer
}

func NewGrpcGateway(logger *slog.Logger, cfg srv.GRPCConfig, ss []srv.ProtoServiceServer) *GrpcGateway {
	return &GrpcGateway{
		logger: logger,
		cfg:    cfg,
		ss:     ss,
	}
}

func (h *GrpcGateway) Handle(mux *http.ServeMux) {
	apiMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true, // use orig names (in snake_case)
		},
	}))
	if len(h.ss) > 0 {
		// TODO context
		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		for _, svc := range h.ss {
			err := svc.RegisterFromEndpoint(ctx, apiMux, h.cfg.Endpoint(), opts)
			if err != nil {
				h.logger.Error("failed to register grpc gateway endpoint", slog.Any("err", err))
			}
		}
	}
	mux.Handle(h.cfg.ApiEndpoint(), apiMux)
}
