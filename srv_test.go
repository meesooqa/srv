package srv

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meesooqa/srv/mocks"
)

func TestHandle_NoMiddleware(t *testing.T) {
	cfg := &mocks.ConfigMock{}
	h1 := &mocks.HandlerMock{}
	h2 := &mocks.HandlerMock{}

	h1.HandleFunc = func(mux *http.ServeMux) {
		mux.HandleFunc("/h1", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("h1"))
		})
	}
	h2.HandleFunc = func(mux *http.ServeMux) {
		mux.HandleFunc("/h2", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("h2"))
		})
	}

	srv := New(cfg, []Handler{h1, h2}, nil)
	h := srv.handle()

	req1 := httptest.NewRequest(http.MethodGet, "/h1", nil)
	res1 := httptest.NewRecorder()
	h.ServeHTTP(res1, req1)
	assert.Equal(t, http.StatusOK, res1.Code)
	assert.Equal(t, "h1", res1.Body.String())

	req2 := httptest.NewRequest(http.MethodGet, "/h2", nil)
	res2 := httptest.NewRecorder()
	h.ServeHTTP(res2, req2)
	assert.Equal(t, http.StatusOK, res2.Code)
	assert.Equal(t, "h2", res2.Body.String())

	assert.Len(t, h1.HandleCalls(), 1)
	assert.Len(t, h2.HandleCalls(), 1)
}

func TestHandle_WithMiddleware(t *testing.T) {
	cfg := &mocks.ConfigMock{}
	h := &mocks.HandlerMock{}
	m1 := &mocks.MiddlewareMock{}
	m2 := &mocks.MiddlewareMock{}

	h.HandleFunc = func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	}

	m1.HandleFunc = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-A", "1")
			next.ServeHTTP(w, r)
		})
	}
	m2.HandleFunc = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-B", "2")
			next.ServeHTTP(w, r)
		})
	}

	srv := New(cfg, []Handler{h}, []Middleware{m1, m2})
	hand := srv.handle()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	hand.ServeHTTP(res, req)

	assert.Equal(t, "1", res.Header().Get("X-A"))
	assert.Equal(t, "2", res.Header().Get("X-B"))
	assert.Equal(t, "ok", res.Body.String())

	assert.Len(t, h.HandleCalls(), 1)
	assert.Len(t, m1.HandleCalls(), 1)
	assert.Len(t, m2.HandleCalls(), 1)
}

func TestRun_Error(t *testing.T) {
	cfg := &mocks.ConfigMock{}
	cfg.HostFunc = func() string { return "invalid_host" }
	cfg.PortFunc = func() int { return 12345 }
	cfg.ReadHeaderTimeoutFunc = func() time.Duration { return time.Second }
	cfg.WriteTimeoutFunc = func() time.Duration { return time.Second }
	cfg.IdleTimeoutFunc = func() time.Duration { return time.Second }

	srv := New(cfg, nil, nil)
	err := srv.Run()
	require.Error(t, err)
}

func TestHandle_NoHandlers_ReturnsServeMux(t *testing.T) {
	cfg := &mocks.ConfigMock{
		HostFunc:              func() string { return "localhost" },
		PortFunc:              func() int { return 8080 },
		ReadHeaderTimeoutFunc: func() time.Duration { return time.Second },
		WriteTimeoutFunc:      func() time.Duration { return time.Second },
		IdleTimeoutFunc:       func() time.Duration { return time.Second },
	}

	srv := New(cfg, nil, nil)
	handler := srv.handle()

	_, ok := handler.(*http.ServeMux)
	assert.True(t, ok)
}

func TestHandle_WithHandler_CallsHandle(t *testing.T) {
	handlerMock := &mocks.HandlerMock{
		HandleFunc: func(mux *http.ServeMux) {},
	}

	cfg := &mocks.ConfigMock{
		HostFunc:              func() string { return "localhost" },
		PortFunc:              func() int { return 8080 },
		ReadHeaderTimeoutFunc: func() time.Duration { return time.Second },
		WriteTimeoutFunc:      func() time.Duration { return time.Second },
		IdleTimeoutFunc:       func() time.Duration { return time.Second },
	}

	srv := New(cfg, []Handler{handlerMock}, nil)
	srv.handle()

	assert.Len(t, handlerMock.HandleCalls(), 1)
	assert.NotNil(t, handlerMock.HandleCalls()[0].Mux)
}

func TestHandle_WithMiddleware_CallsMiddleware(t *testing.T) {
	middlewareMock := &mocks.MiddlewareMock{
		HandleFunc: func(next http.Handler) http.Handler {
			return next
		},
	}

	cfg := &mocks.ConfigMock{
		HostFunc:              func() string { return "localhost" },
		PortFunc:              func() int { return 8080 },
		ReadHeaderTimeoutFunc: func() time.Duration { return time.Second },
		WriteTimeoutFunc:      func() time.Duration { return time.Second },
		IdleTimeoutFunc:       func() time.Duration { return time.Second },
	}

	srv := New(cfg, nil, []Middleware{middlewareMock})
	srv.handle()

	calls := middlewareMock.HandleCalls()
	assert.Len(t, calls, 1)

	next, ok := calls[0].Next.(*http.ServeMux)
	assert.True(t, ok)
	assert.NotNil(t, next)
}

func TestHandle_WithMultipleMiddlewares_CallsInOrder(t *testing.T) {
	mw1 := &mocks.MiddlewareMock{
		HandleFunc: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
	}

	mw2 := &mocks.MiddlewareMock{
		HandleFunc: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
	}

	cfg := &mocks.ConfigMock{
		HostFunc:              func() string { return "localhost" },
		PortFunc:              func() int { return 8080 },
		ReadHeaderTimeoutFunc: func() time.Duration { return time.Second },
		WriteTimeoutFunc:      func() time.Duration { return time.Second },
		IdleTimeoutFunc:       func() time.Duration { return time.Second },
	}

	srv := New(cfg, nil, []Middleware{mw1, mw2})
	srv.handle()

	assert.Len(t, mw1.HandleCalls(), 1)
	assert.Len(t, mw2.HandleCalls(), 1)
}

func TestRun_CallsConfigMethods(t *testing.T) {
	cfg := &mocks.ConfigMock{
		HostFunc:              func() string { return "localhost" },
		PortFunc:              func() int { return 0 },
		ReadHeaderTimeoutFunc: func() time.Duration { return time.Second },
		WriteTimeoutFunc:      func() time.Duration { return time.Second },
		IdleTimeoutFunc:       func() time.Duration { return time.Second },
	}

	srv := New(cfg, nil, nil)

	errCh := make(chan error)
	go func() {
		errCh <- srv.Run()
	}()

	time.Sleep(100 * time.Millisecond)

	assert.Len(t, cfg.HostCalls(), 1)
	assert.Len(t, cfg.PortCalls(), 1)
	assert.Len(t, cfg.ReadHeaderTimeoutCalls(), 1)
	assert.Len(t, cfg.WriteTimeoutCalls(), 1)
	assert.Len(t, cfg.IdleTimeoutCalls(), 1)
}

func TestRun_ReturnsErrorOnInvalidHost(t *testing.T) {
	cfg := &mocks.ConfigMock{
		HostFunc:              func() string { return "invalid_host" },
		PortFunc:              func() int { return 8080 },
		ReadHeaderTimeoutFunc: func() time.Duration { return time.Second },
		WriteTimeoutFunc:      func() time.Duration { return time.Second },
		IdleTimeoutFunc:       func() time.Duration { return time.Second },
	}

	srv := New(cfg, nil, nil)

	err := srv.Run()
	assert.Error(t, err)
}
