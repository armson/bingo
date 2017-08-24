package redis

import (
	"github.com/armson/bingo/utils"
	"time"
)

func (client *Redis) Get(key string) (string, error) {
	cmd := client.pool(key).Get(key)
	client.logs(cmd.String())
	return cmd.Result()
}

func (client *Redis) Set(key string, value interface{}) (string, error) {
	cmd := client.pool(key).Set(key, value, 0)
	client.logs(cmd.String())
	return cmd.Result()
}

func (client *Redis) SetEx(key string, value interface{}, seconds int) (string, error) {
	timeString := utils.String.Join(utils.Int.String(seconds), "s")
	duration := utils.Duration.Parse(timeString)
	cmd := client.pool(key).Set(key, value, duration)
	client.logs(cmd.String())
	return cmd.Result()
}

// example
// SetNx("k","val","3600s")
// SetNx("k",1)
func (client *Redis) SetNX(key string, args ...interface{}) (bool, error) {
	if len(args) > 1 {
		duration, _ := time.ParseDuration(args[1].(string))
		cmd := client.pool(key).SetNX(key, args[0], duration)
		client.logs(cmd.String())
		return cmd.Result()
	} else {
		cmd := client.pool(key).SetNX(key, args[0], 0)
		client.logs(cmd.String())
		return cmd.Result()
	}
}

func (client *Redis) SMembers(key string) ([]string, error) {
	cmd := client.pool(key).SMembers(key)
	client.logs(cmd.String())
	return cmd.Result()
}

func (client *Redis) SAdd(key string, members ...interface{}) (int64, error) {
	cmd := client.pool(key).SAdd(key, members...)
	client.logs(cmd.String())
	return cmd.Result()
}

func (client *Redis) Del(key string) (int64, error) {
	cmd := client.pool(key).Del(key)
	client.logs(cmd.String())
	return cmd.Result()
}

func (client *Redis) Expire(key, expiration string) (bool, error) {
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return false, err
	}
	cmd := client.pool(key).Expire(key, duration)
	client.logs(cmd.String())
	return cmd.Result()
}

func (client *Redis) HGet(key string, mobile string) (string, error) {
	cmd := client.pool(key).HGet(key, mobile)
	client.logs(cmd.String())
	return cmd.Result()
}
