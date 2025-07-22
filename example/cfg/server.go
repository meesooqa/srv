package cfg

import "time"

// Server contains server configuration
type Server struct {
	RawHost              string        `yaml:"host"`
	RawPort              int           `yaml:"port"`
	RawReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	RawWriteTimeout      time.Duration `yaml:"write_timeout"`
	RawIdleTimeout       time.Duration `yaml:"idle_timeout"`
}

func (cfg *Server) Host() string {
	return cfg.RawHost
}

func (cfg *Server) Port() int {
	return cfg.RawPort
}

func (cfg *Server) ReadHeaderTimeout() time.Duration {
	return cfg.RawReadHeaderTimeout
}

func (cfg *Server) WriteTimeout() time.Duration {
	return cfg.RawWriteTimeout
}

func (cfg *Server) IdleTimeout() time.Duration {
	return cfg.RawIdleTimeout
}
