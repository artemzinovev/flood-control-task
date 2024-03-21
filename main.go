package main

import (
	"context"
	"fmt"
	floodcontrol "task/floodControl"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	var user int64
	user = 5
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:8080",
		Password: "",
		DB:       0,
	})

	checkFlood := floodcontrol.CheckFloodConstructor(30*time.Second, 10, redis)

	check, err := checkFlood.Check(context.Background(), user)
	if err != nil {
		fmt.Printf("Error check: %s", err)
		return
	}

	if check == true {
		fmt.Printf("User: %v allowed", user)
	}
	if check == false {
		fmt.Printf("User: %v banned", user)
	}

}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
