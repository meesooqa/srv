package middlewares

import "net/http"

type CORSConfig interface {
	AllowedOrigins() []string
}

// CORS middleware to handle Cross-Origin Resource Sharing
type CORS struct {
	allowedOrigins []string
}

func NewCORS(cfg CORSConfig) *CORS {
	return &CORS{
		allowedOrigins: cfg.AllowedOrigins(),
	}
}

func (m *CORS) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(m.allowedOrigins) <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")
		isAllowed := false
		for _, allowed := range m.allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}
		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Accept, Content-Type, Content-Length, Accept-Encoding, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
