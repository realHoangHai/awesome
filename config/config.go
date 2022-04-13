package config

import (
	"bufio"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
)

var (
	cfg = &Config{}
)

// Read read configuration from environment to the target ptr.
func Read(ptr interface{}, opts ...ReadOption) error {
	return cfg.Read(ptr, opts...)
}

// Close close the default config reader.
func Close() error {
	return cfg.Close()
}

type Config struct {
}

// Read implements config.Reader interface.
func (c *Config) Read(ptr interface{}, opts ...ReadOption) error {
	ops := &ReadOptions{}
	ops.Apply(opts...)
	if ops.File != "" {
		if err := loadEnvFromFile(ops.File); err != nil {
			return err
		}
	}
	if ops.FileNoErr != "" {
		_ = loadEnvFromFile(ops.FileNoErr)
	}
	return envconfig.Process(ops.Prefix, ptr)
}

// loadEnvFromFile load environments from file
// and set them to system environment via os.Setenv.
func loadEnvFromFile(f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "#") || strings.TrimSpace(txt) == "" {
			continue
		}
		env := strings.SplitN(txt, "=", 2)
		if len(env) != 2 {
			return fmt.Errorf("invalid pair: %v", txt)
		}
		k := env[0]
		v := env[1]
		_ = os.Setenv(k, v)
	}
	return nil
}

// Close implements config.Reader interface.
func (c *Config) Close() error {
	return nil
}
