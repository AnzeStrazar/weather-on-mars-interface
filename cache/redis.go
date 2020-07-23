package cache

import (
	"log"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

// Constructor for the Redis subscription.
func NewRedis(redisEndpoint string) *Redis {

	rds := &Redis{}

	options, err := redis.ParseURL(redisEndpoint)

	if err != nil {
		log.Panic("Parsing redis URL had failed.")
	}

	rds.client = redis.NewClient(options)

	// Check that redis endpoint is in fact up and responsive
	// by pinging it.
	ping := rds.client.Ping()
	_, err = ping.Result()

	if err != nil {
		log.Panic("Could not establish connection to the redis.")
	} else {
		log.Println("Connection to redis established successfully.")
	}

	return rds
}

// Returns string representation of the key stored in redis cache
func (r *Redis) Get(key string) string {
	return r.client.Get(key).Val()
}

// Sets the key to the redis cache.
func (r *Redis) Set(key string, value interface{}) {
	r.client.Set(key, value, 0)
}
