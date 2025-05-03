package config

// Config holds the server configuration parameters.
type Config struct {
	Host     string // Host address to bind to
	GRPCPort string // Port for gRPC server
	HTTPPort string // Port for HTTP gateway server
}

// NewConfig creates a new Config instance with the given parameters.
func NewConfig(host, grpcPort, httpPort string) *Config {
	return &Config{
		Host:     host,
		GRPCPort: grpcPort,
		HTTPPort: httpPort,
	}
} 