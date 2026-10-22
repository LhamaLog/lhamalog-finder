package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	Pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	ultimoID := "110"

	lista, err := rdb.XRead(context.Background(), &redis.XReadArgs{
		Streams: []string{"fila:logs", ultimoID},
		Block:   0,
		Count:   10,
	}).Result()

	if err != nil {
		log.Fatalf("Failed to read from Redis: %v", err)
	}

	fmt.Println(lista)

	fmt.Println("================================================")

	for _, lista := range lista {
		for _, message := range lista.Messages {
			fmt.Println(message.ID)
			fmt.Println(message.Values["id"])

			id := message.ID

			errs := rdb.XDel(context.Background(), "fila:logs", id).Err()
			if errs != nil {
				log.Fatalf("Failed to delete from Redis: %v", errs)
			}
		}
	}

	fmt.Println(Pong)
	return rdb
}
