// worker.go

package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os/exec"
)

var redisClientWorker *redis.Client

func runWorker() {
	for {
		result, err := redisClientWorker.BLPop(context.Background(), 0, "image_queue").Result()
		if err != nil {
			fmt.Println("Error getting item from queue:", err)
			continue
		}

		imageName := result[1]

		cmd := exec.CommandContext(context.Background(), "bash", "run.sh", imageName)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error running command:", err)
			continue
		}

		fmt.Printf("Image %s executed successfully!\n", imageName)
	}
}

func main() {
	// Kết nối đến Redis
	redisClientWorker = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Chạy worker
	runWorker()
}
