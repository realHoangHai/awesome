package main

import (
	"context"
	"fmt"
	"github.com/realHoangHai/awesome/config"
)

func main() {
	cfg := &config.Config{
		Core:  config.SectionCore{},
		API:   config.SectionAPI{},
		DB:    config.SectionDB{},
		Log:   config.SectionLog{},
		Redis: config.SectionRedis{},
	}
	err := InitializeServer(context.Background(), cfg)
	fmt.Println(err)
}
