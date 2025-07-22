package middlewares

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

type Logging struct {
	logger *slog.Logger
}

func NewLogging(logger *slog.Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (m *Logging) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("received request",
			slog.String("method", r.Method),
			slog.String("URL.Path", r.URL.Path))
		// body
		if r.Method == "PATCH" || r.Method == "POST" || r.Method == "PUT" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				m.logger.Error("failed to read request body", slog.String("error", err.Error()))
			} else {
				// restore body
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				m.logger.Info("request body", slog.String("body", string(bodyBytes)))
			}
		}
		next.ServeHTTP(w, r)
	})
}
