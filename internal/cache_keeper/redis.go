package cache_keeper

import (
	"github.com/go-redis/redis/v7"
	"github.com/golang/protobuf/proto"
)

type redisCacheKeeper struct {
	cluster *redis.ClusterClient
}

func NewRedisCacheKeeper(opts *redis.ClusterOptions) CacheKeeper {
	cluster :=  redis.NewClusterClient(opts)
	return &redisCacheKeeper{cluster:cluster}
}

func (r *redisCacheKeeper) Get(key string, msg proto.Message) error {
	 cmd := r.cluster.Do("GET", key)
	 if cmd.Err() != nil {
	 	if cmd.Err() == redis.Nil {
	 		return &cacheError{isNotFound:true}
		}
	 	return cmd.Err()
	 }
	 reply, err := cmd.String()
	 if err != nil {
	 	return err
	 }

	 err = proto.Unmarshal([]byte(reply), msg)
	 if err != nil {
	 	return err
	 }

	 return nil
}

func (r *redisCacheKeeper) Set(key string, message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	cmd := r.cluster.Do("SETEX", key, 10, string(data))
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return err
	}

	return nil
}