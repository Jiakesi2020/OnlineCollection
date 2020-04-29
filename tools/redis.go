package tools

import "github.com/go-redis/redis"

func GetRedisInstance(hosts []string) (*redis.Client, error) {
	cc := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    hosts[0],
	})

	if err := cc.Ping().Err(); err != nil {
		return nil, err
	}

	return cc, nil
}
