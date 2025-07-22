package cfg

// GRPCServer contains gRPC server configuration
type GRPCServer struct {
	RawEndpoint    string `yaml:"endpoint"`
	RawApiEndpoint string `yaml:"api_endpoint"`
}

func (cfg *GRPCServer) Endpoint() string {
	return cfg.RawEndpoint
}

func (cfg *GRPCServer) ApiEndpoint() string {
	return cfg.RawApiEndpoint
}
