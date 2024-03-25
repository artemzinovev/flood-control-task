package main

import (
	"context"
	"log/slog"
	"os"
	"task/internal/config"
	redisRepository "task/internal/repository/redis"
	"task/internal/service"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := config.MustLoad(log)

	RequestCountRepository, err := redisRepository.NewRequestCountRepository(
		cfg.RedisConfig.Addr,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
		cfg.RedisConfig.TTL,
	)
	if err != nil {
		log.Error("failed to create redis repository", err)
		os.Exit(1)
	}

	floodControl := service.NewFloodControlService(RequestCountRepository, cfg.RequestLimit)

	for {
		go func() {
			ok, err := floodControl.Check(context.Background(), 3)
			if err != nil {
				log.Error("error", err)
				os.Exit(1)
			}

			if !ok {
				log.Info("flood control failed")
				os.Exit(1)
			}

			log.Info("flood control passed")
		}()
	}

}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
