package srv

import (
	"fmt"
	"net/http"
	"time"
)

//go:generate moq --out ./mocks/config_mock.go --pkg mocks --skip-ensure --with-resets -fmt goimports . Config
//go:generate moq --out ./mocks/handler_mock.go --pkg mocks --skip-ensure --with-resets -fmt goimports . Handler
//go:generate moq --out ./mocks/middleware_mock.go --pkg mocks --skip-ensure --with-resets -fmt goimports . Middleware

type Config interface {
	Host() string
	Port() int

	ReadHeaderTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
}

type Handler interface {
	Handle(mux *http.ServeMux)
}

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type Server struct {
	cfg Config
	hh  []Handler
	mw  []Middleware
}

func New(cfg Config, hh []Handler, mw []Middleware) *Server {
	return &Server{
		cfg: cfg,
		hh:  hh,
		mw:  mw,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.cfg.Host(), s.cfg.Port()),
		Handler:           s.handle(),
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout(),
		WriteTimeout:      s.cfg.WriteTimeout(),
		IdleTimeout:       s.cfg.IdleTimeout(),
	}
	return srv.ListenAndServe()
}

func (s *Server) handle() http.Handler {
	mux := http.NewServeMux()
	for _, handler := range s.hh {
		handler.Handle(mux)
	}

	if len(s.mw) > 0 {
		middleHandler := http.Handler(mux)
		for _, middleware := range s.mw {
			middleHandler = middleware.Handle(middleHandler)
		}
		return middleHandler
	} else {
		return mux
	}
}
