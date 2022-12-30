package switchbot

import (
	"context"
	"sync"
	"time"
)

type cache struct {
	mu sync.RWMutex

	c Client

	ttl time.Duration

	expire time.Time
	temp   float32
	hum    float32
}

var _ Client = (*cache)(nil)

func NewCacheClient(c Client) Client {
	return &cache{
		c:      c,
		ttl:    60 * time.Second,
		expire: time.Now(),
		temp:   0.0,
		hum:    0.0,
	}
}

func (c *cache) GetMetrics(ctx context.Context) (float32, float32, error) {
	if c.expire.After(time.Now()) {
		t, h := c.get()
		return t, h, nil
	}

	temp, hum, err := c.c.GetMetrics(ctx)

	c.set(temp, hum)

	c.expire = time.Now().Add(c.ttl)

	return temp, hum, err
}

func (c *cache) get() (float32, float32) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.temp, c.hum
}

func (c *cache) set(temp, hum float32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.temp = temp
	c.hum = hum
}
