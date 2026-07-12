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

	metrics *Metrics
}

var _ Client = (*cache)(nil)

func NewCacheClient(c Client) Client {
	return &cache{
		c:       c,
		ttl:     60 * time.Second,
		expire:  time.Now(),
		metrics: nil,
	}
}

func (c *cache) GetMetrics(ctx context.Context) (*Metrics, error) {
	if c.isFresh() {
		m := c.get()
		return m, nil
	}

	m, err := c.c.GetMetrics(ctx)
	if err != nil {
		return nil, err
	}

	c.setMetrics(m, time.Now().Add(c.ttl))

	return &Metrics{
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
	}, err
}

func (c *cache) isFresh() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.expire.After(time.Now())
}

func (c *cache) get() *Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics
}

func (c *cache) setMetrics(m *Metrics, expire time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = m
	c.expire = expire
}
