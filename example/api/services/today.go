package services

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/meesooqa/srv/example/api/gen/pb/today/v1"
)

// TodaySS implements the gRPC server service for date operations
type TodaySS struct {
	pb.UnimplementedTodayServiceServer
}

// NewTodaySS creates a new TodaySS client
func NewTodaySS() *TodaySS {
	return &TodaySS{}
}

func (s *TodaySS) Register(grpcServer *grpc.Server) {
	pb.RegisterTodayServiceServer(grpcServer, s)
}

func (*TodaySS) RegisterFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return pb.RegisterTodayServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

func (*TodaySS) Today(_ context.Context, r *pb.TodayRequest) (*pb.TodayResponse, error) {
	layout := r.Format
	if layout == "" {
		layout = "2006-01-02"
	}
	today := time.Now().Format(layout)

	return &pb.TodayResponse{
		Today:  today,
		Format: layout,
	}, nil
}
