package main

import (
	"github.com/Ryan-Har/chaos-kube/chaos-executor/services"
	"github.com/Ryan-Har/chaos-kube/pkg/common"
	"github.com/Ryan-Har/chaos-kube/pkg/config"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.Load("executor")
	if err != nil {
		slog.Error("unable to load config", "err", err)
		os.Exit(1)
	}

	redisClient := common.NewRedisClient(cfg.RedisConfig)
	executor := services.NewExecutorService(cfg, redisClient)
	executor.Start()
}
