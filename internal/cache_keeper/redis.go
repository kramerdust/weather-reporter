package cache_keeper

import (
	"github.com/chasex/redis-go-cluster"
	"github.com/golang/protobuf/proto"
)

type redisCacheKeeper struct {
	cluster redis.Cluster
}

func NewRedisCacheKeeper(opts *redis.Options) (CacheKeeper, error) {
	cluster, err :=  redis.NewCluster(opts)
	if err != nil {
		return nil, err
	}

	return &redisCacheKeeper{cluster:cluster}, nil
}

func (r *redisCacheKeeper) Get(key string, msg proto.Message) error {
	 reply, err := redis.Bytes(r.cluster.Do("GET", key))
	 if err != nil {
	 	return err
	 }

	 err = proto.Unmarshal(reply, msg)
	 if err != nil {
	 	return err
	 }
}

func (r *redisCacheKeeper) Set(key string, message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	_, err = r.cluster.Do("SET", key, data)
	if err != nil {
		return err
	}
}