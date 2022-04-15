package main

import (
	"fmt"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/internal/repo"
	"github.com/realHoangHai/awesome/internal/server"
	"github.com/realHoangHai/awesome/internal/service"
	"github.com/realHoangHai/awesome/pkg/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	var services []server.Service
	entClient := repo.NewEntClient(&cfg)
	redisClient := repo.NewRedisCmd(&cfg)
	store, closeFunc, err := repo.NewStore(entClient, redisClient)
	defer closeFunc()
	userStore := repo.NewUserRepo(store)
	cardStore := repo.NewCardRepo(store)
	addressStore := repo.NewAddressRepo(store)

	userBiz := biz.NewUserBiz(userStore)
	cardBiz := biz.NewCardBiz(cardStore)
	addressBiz := biz.NewAddressBiz(addressStore)

	services = append(services, service.NewUserService(userBiz, addressBiz, cardBiz))

	s := server.New(server.FromEnv(&cfg))
	if err := s.ListenAndServe(services...); err != nil {
		log.Fatal(err)
	}

	fmt.Println(err)
}
