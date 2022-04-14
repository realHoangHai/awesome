package config

import "time"

type Config struct {
	Core  SectionCore  `mapstructure:"core"`
	API   SectionAPI   `mapstructure:"api"`
	DB    SectionDB    `mapstructure:"db"`
	Log   SectionLog   `mapstructure:"log"`
	Redis SectionRedis `mapstructure:"redis"`
}

type SectionCore struct {
	// Name is name of the service.
	Name string `mapstructure:"NAME" default:"awesome"`
	// Address is the address of the service in form of host:port.
	// If PORT environment variable is configured, it will be prioritized over ADDRESS.
	Address string `mapstructure:"ADDRESS" default:":8000"`
	// TLSCertFile is path to the TLS certificate file.
	TLSCertFile string `mapstructure:"TLS_CERT_FILE"`
	// TLSKeyFile is the path to the TLS key file.
	TLSKeyFile string `mapstructure:"TLS_KEY_FILE"`

	// ReadTimeout is read timeout of both gRPC and HTTP server.
	ReadTimeout time.Duration `mapstructure:"READ_TIMEOUT" default:"30s"`
	// WriteTimeout is write timeout of both gRPC and HTTP server.
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT" default:"30s"`
	//ShutdownTimeout is timeout for shutting down the server.
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT" default:"30s"`
	// APIPrefix is path prefix that gRPC API Gateway is routed to.
	APIPrefix string `envconfig:"API_PREFIX" default:"/api/"`

	// JWTSecret is a short way to enable JWT Authentictor with the secret.
	JWTSecret string `mapstructure:"JWT_SECRET"`
	// ContextLogger is an option to enable context logger with request-id.
	ContextLogger bool `mapstructure:"CONTEXT_LOGGER" default:"true"`

	// Recovery is a short way to enable recovery interceptors for both unary and stream handlers.
	Recovery bool `mapstructure:"RECOVERY" default:"true"`

	// CORS options
	CORSAllowedHeaders    []string `mapstructure:"CORS_ALLOWED_HEADERS"`
	CORSAllowedMethods    []string `mapstructure:"CORS_ALLOWED_METHODS"`
	CORSAllowedOrigins    []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedCredential bool     `mapstructure:"CORS_ALLOWED_CREDENTIAL" default:"false"`

	// Metrics enable/disable standard metrics
	Metrics bool `mapstructure:"METRICS" default:"true"`
	// MetricsPath is API path for Prometheus metrics.
	MetricsPath string `mapstructure:"METRICS_PATH" default:"/internal/metrics"`

	RoutesPrioritization bool   `mapstructure:"ROUTES_PRIORITIZATION" default:"true"`
	ShutdownHook         string `mapstructure:"SHUTDOWN_HOOK"`
}

type SectionAutoTLS struct {
	Enabled bool   `mapstructure:"enabled"`
	Folder  string `mapstructure:"folder"`
	Host    string `mapstructure:"host"`
}

type SectionAPI struct {
}

type SectionDB struct {
	Driver string `mapstructure:"driver"`
	Source string `mapstructure:"source"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Format string `mapstructure:"format"`
}

type SectionRedis struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}
