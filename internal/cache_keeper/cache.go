package cache_keeper

import "github.com/golang/protobuf/proto"


type CacheKeeper interface {
	Get(key string, message proto.Message) error
	Set(key string, message proto.Message) error
}

type cacheError struct {
	internal error
	isNotFound bool
}

func (ce *cacheError) Error() string {
	if ce.isNotFound {
		return "key not found"
	}
	if ce.internal != nil {
		return ce.internal.Error()
	}
	return ""
}

func CacheNotFound(err error) bool {
	if ce, ok := err.(*cacheError); ok {
		return ce.isNotFound
	}
	return false
}