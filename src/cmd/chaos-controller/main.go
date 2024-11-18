package main

import (
	"github.com/Ryan-Har/chaos-kube/chaos-controller/handler"
	"github.com/Ryan-Har/chaos-kube/chaos-controller/services"
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.Load("controller")
	if err != nil {
		slog.Error("unable to load config", "err", err)
		os.Exit(1)
	}

	responseManager := common.NewResponseManager()
	redisClient := common.NewRedisClient(cfg.RedisConfig)
	hndlr := handler.NewChaosHandler(redisClient, responseManager)
	controller := services.NewControllerService(cfg, hndlr)
	controller.Start()
}
