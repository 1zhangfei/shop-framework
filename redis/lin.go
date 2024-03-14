package redis

import (
	"2108a-zg5/week3/shop/framework/config"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func WithRedisInit(address string, hand func(cli *redis.Client) (string, error)) (string, error) {
	if err := config.ViperInit(address); err != nil {
		return "", err
	}

	var app struct {
		RedisConf struct {
			Host string
			Port string
		} `json:"Redis"`
	}

	id := viper.GetString("Database.DataId")
	Group := viper.GetString("Database.Group")
	res, err := config.GetConfig(id, Group)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(res), &app); err != nil {
		return "", err
	}

	Cli := redis.NewClient(&redis.Options{
		Addr: app.RedisConf.Host + ":" + app.RedisConf.Port,
		DB:   6,
	})
	defer Cli.Close()

	res, err = hand(Cli)
	if err != nil {
		return "", err
	}
	return res, nil
}
