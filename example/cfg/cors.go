package cfg

// CORS contains CORS configuration
type CORS struct {
	RawAllowedOrigins []string `yaml:"allowed_origins"`
}

func (cfg *CORS) AllowedOrigins() []string {
	return cfg.RawAllowedOrigins
}
