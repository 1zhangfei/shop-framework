package redis

import (
	"context"
	"encoding/json"
	"github.com/1zhangfei/shop-framework/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

func WithRedisInit(address string, hand func(cli *redis.Client) error) error {
	if err := config.ViperInit(address); err != nil {
		return err
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
		return err
	}

	if err = json.Unmarshal([]byte(res), &app); err != nil {
		return err
	}

	Cli := redis.NewClient(&redis.Options{
		Addr: app.RedisConf.Host + ":" + app.RedisConf.Port,
		DB:   6,
	})
	defer Cli.Close()

	err = hand(Cli)
	if err != nil {
		return err
	}
	return nil
}

func ExisKey(ctx context.Context, Address, Key string) (bool, error) {
	var data int64
	var err error

	err = WithRedisInit(Address, func(cli *redis.Client) error {
		data, err = cli.Exists(ctx, Key).Result()
		return err
	})
	if err != nil {
		return false, err
	}
	if data > 0 {
		return true, nil
	}
	return false, nil
}

func GetByVal(ctx context.Context, Address, key string) (string, error) {
	var data string
	var err error
	err = WithRedisInit(Address, func(cli *redis.Client) error {
		data, err = cli.Get(ctx, key).Result()
		return err
	})

	return data, err
}

func SetKey(ctx context.Context, Address, key string, Val interface{}, duration time.Duration) error {
	return WithRedisInit(Address, func(cli *redis.Client) error {
		return cli.Set(ctx, key, Val, duration).Err()
	})
}

func Lock(ctx context.Context, Address, key string, val interface{}, duration time.Duration, isContinue bool) (bool, error) {
	var re = false
	err := WithRedisInit(Address, func(cli *redis.Client) error {
		if !isContinue {
			for {
				res, err := cli.SetNX(ctx, key, val, duration).Result()
				if err != nil {
					return err
				}
				re = true
				if res {
					return nil
				}
			}
		}
		res, err := cli.SetNX(ctx, key, val, duration).Result()
		re = res
		return err
	})
	return re, err
}

func UnLock(ctx context.Context, Address, key string) error {
	return WithRedisInit(Address, func(cli *redis.Client) error {
		return cli.Del(ctx, key).Err()
	})
}
