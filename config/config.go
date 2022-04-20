package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Core   SectionCore   `mapstructure:"core"`
	API    SectionAPI    `mapstructure:"api"`
	DB     SectionDB     `mapstructure:"db"`
	Log    SectionLog    `mapstructure:"log"`
	Redis  SectionRedis  `mapstructure:"redis"`
	Health SectionHealth `mapstructure:"health"`
}

func LoadConfig(path string) (cfg Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}

type SectionCore struct {
	Name        string `mapstructure:"name"`
	Address     string `mapstructure:"address"`
	TLSCertFile string `mapstructure:"tls_cert_file"`
	TLSKeyFile  string `mapstructure:"tls_key_file"`

	HealthCheckPath string `mapstructure:"health_check_path"`

	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	APIPrefix       string        `mapstructure:"api_prefix"`

	WebDir    string `mapstructure:"web_dir"`
	WebIndex  string `mapstructure:"web_index"`
	WebPrefix string `mapstructure:"web_prefix"`

	JWTSecret     string `mapstructure:"jwt_secret"`
	ContextLogger bool   `mapstructure:"context_logger"`

	Recovery bool `mapstructure:"recovery"`

	CORSAllowedHeaders    []string `mapstructure:"cors_allowed_headers"`
	CORSAllowedMethods    []string `mapstructure:"cor_allowed_methods"`
	CORSAllowedOrigins    []string `mapstructure:"cors_allowed_origins"`
	CORSAllowedCredential bool     `mapstructure:"cors_allowed_credential"`

	PProf       bool   `mapstructure:"pprof"`
	PProfPrefix string `mapstructure:"pprof_prefix"`

	Metrics     bool   `mapstructure:"metrics"`
	MetricsPath string `mapstructure:"metrics_path" `

	RoutesPrioritization bool   `mapstructure:"routes_prioritization" `
	ShutdownHook         string `mapstructure:"shutdown_hook"`
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
	Level      int               `mapstructure:"log_level"`
	Format     string            `mapstructure:"log_format"`
	TimeFormat string            `mapstructure:"log_time_format"`
	Output     string            `mapstructure:"log_output"`
	Fields     map[string]string `mapstructure:"log_fields"`
}

type SectionRedis struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type SectionHealth struct {
	Interval time.Duration `mapstructrue:"interval"`
	Timeout  time.Duration `mapstructrue:"timeout"`
}
