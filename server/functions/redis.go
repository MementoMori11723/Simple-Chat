package functions

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func connect() *redis.Client {
  rdb := redis.NewClient(&redis.Options{
    Addr: "simple_chat-redis:6379",
    Password: "",
    DB: 0,
  })
  return rdb
}

func TestConnection() {
  rdb := connect()
  defer rdb.Close()
  pong, err := rdb.Ping(context.Background()).Result()
  if err != nil {
    slog.Error("Error in connecting to redis", "Error - msg", err.Error())
    os.Exit(1)
  } else {
    slog.Info("Connected to redis", "Ping", pong)
  }
}

func Add(key string, value string, expire time.Duration) {
  rdb := connect()
  err := rdb.Set(context.Background(), key, value, expire).Err()
  if err != nil {
    slog.Error("Error in setting key value in redis", "Error - msg", err.Error())
    os.Exit(1)
  } else {
    slog.Info("Key value set in redis", "Key", key)
  }
}

func Get_data(key string) string {
  rdb := connect()
  val, err := rdb.Get(context.Background(), key).Result()
  if err != nil {
    slog.Error("Error in getting value from redis", "Error - msg", err.Error())
    return ""
  } else {
    slog.Info("Value fetched from redis", "Key", key)
  }
  return val
}

func Get_keys(key string) string {
  rdb := connect()
  val, err := rdb.Keys(context.Background(), key).Result()
  if err != nil {
    slog.Error("Error in getting value from redis", "Error - msg", err.Error())
    return ""
  }
  if len(val) > 1 {
    slog.Error("Error in getting value from redis", "Error - msg", "Too many Values!")
    return ""
  }
  return val[0]
}

func Delete(key string) {
  rdb := connect()
  err := rdb.Del(context.Background(), key).Err()
  if err != nil {
    slog.Error("Error in deleting key from redis", "Error - msg", err.Error())
    os.Exit(1)
  } else {
    slog.Info("Key deleted from redis", "Key", key)
  }
}
