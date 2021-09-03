package cache

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

// Memory 临时用于兼容 github.com/NICEXAI/WeChatCustomerServiceSDK/cache
type Memory struct {
	ins *cache.Cache
}

func New() *Memory {
	return &Memory{
		ins: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *Memory) Set(k, v string, expires time.Duration) error {
	c.ins.Set(k, v, expires)
	return nil
}

func (c *Memory) Get(k string) (string, error) {
	v, ok := c.ins.Get(k)
	if !ok {
		return "", nil
	}
	if vStr, ok := v.(string); ok {
		return vStr, nil
	}
	return "", errors.New("cache value type error")
}
