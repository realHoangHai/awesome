package config_test

import (
	"github.com/realHoangHai/awesome/config"
	"time"
)

func Example() {
	var conf struct {
		Name        string            `envconfig:"NAME" default:"awesome"`
		Address     string            `envconfig:"ADDRESS" default:"0.0.0.0:8088"`
		Secret      string            `envconfig:"SECRET"`
		Fields      []string          `envconfig:"FIELDS" default:"field1,field2"`
		ReadTimeout time.Duration     `envconfig:"READ_TIMEOUT" default:"30s"`
		Enable      bool              `envconfig:"ENABLE" default:"true"`
		Map         map[string]string `envconfig:"MAP" default:"key:value,key1:value1"`
	}
	_ = config.Read(&conf)
}

func ExampleRead_withOptions() {
	var conf struct {
		Name    string `envconfig:"NAME" default:"awesome"`
		Address string `envconfig:"ADDRESS" default:"0.0.0.0:8088"`
		Secret  string `envconfig:"SECRET"`
	}
	_ = config.Read(&conf, config.WithPrefix("HTTP"))
}
