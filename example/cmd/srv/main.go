package main

import (
	"log"
	"log/slog"

	"github.com/meesooqa/lgr"
	"github.com/meesooqa/srv"
	"github.com/meesooqa/srv/handlers"
	"github.com/meesooqa/srv/middlewares"

	"github.com/meesooqa/srv/example/api/services"
	"github.com/meesooqa/srv/example/cfg"
)

func main() {
	conf, err := cfg.Load("etc/config.yml")
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}
	logger, closer := lgr.New(conf.Log)
	if closer != nil {
		defer func() { _ = closer.Close() }()
	}

	ss := []srv.ProtoServiceServer{
		services.NewTodaySS(),
	}
	grpcSrv := srv.NewGRPCServer(logger, conf.GRPCServer, ss)
	err = grpcSrv.Run()
	if err != nil {
		log.Fatal(err)
	}

	hh := []srv.Handler{
		handlers.NewGrpcGateway(logger, conf.GRPCServer, ss),
	}
	mw := []srv.Middleware{
		middlewares.NewLogging(logger),
		middlewares.NewCORS(conf.CORS),
	}
	s := srv.New(conf.Server, hh, mw)

	logger.Info("server started", slog.String("host", conf.Server.Host()), slog.Int("port", conf.Server.Port()))
	log.Fatal(s.Run())
}
