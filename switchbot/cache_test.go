package switchbot

import (
	"context"
	"errors"
	"testing"
)

type clientFunc func(context.Context) (*Metrics, error)

func (f clientFunc) GetMetrics(ctx context.Context) (*Metrics, error) {
	return f(ctx)
}

func TestCacheStoresSuccessfulMetrics(t *testing.T) {
	calls := 0
	client := clientFunc(func(context.Context) (*Metrics, error) {
		calls++
		if calls == 1 {
			return &Metrics{23.5, 48}, nil
		}

		return &Metrics{99, 99}, nil
	})
	cache := NewCacheClient(client)

	m, err := cache.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected fresh metrics, got %v", err)
	}
	if m.Temperature != 23.5 || m.Humidity != 48 {
		t.Fatalf("expected fresh metrics, got m.Temperature=%v hum=%v", m.Temperature, m.Humidity)
	}

	m, err = cache.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected cached metrics, got %v", err)
	}
	if m.Temperature != 23.5 || m.Humidity != 48 {
		t.Fatalf("expected cached metrics, got temp=%v hum=%v", m.Temperature, m.Humidity)
	}
	if calls != 1 {
		t.Fatalf("expected successful metrics to be cached, got %d calls", calls)
	}
}

func TestCacheDoesNotStoreFailedMetrics(t *testing.T) {
	temporaryErr := errors.New("temporary error")
	calls := 0
	client := clientFunc(func(context.Context) (*Metrics, error) {
		calls++
		if calls == 1 {
			return nil, temporaryErr
		}

		return &Metrics{24.5, 51}, nil
	})
	cache := NewCacheClient(client)

	m, err := cache.GetMetrics(context.Background())
	if !errors.Is(err, temporaryErr) {
		t.Fatalf("expected temporary error, got %v", err)
	}
	if m != nil {
		t.Fatalf(
			"failed response should return zero values, got temp=%v hum=%v",
			m.Temperature, m.Humidity,
		)
	}

	m, err = cache.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	if m.Temperature != 24.5 || m.Humidity != 51 {
		t.Fatalf(
			"expected fresh metrics after failure, got temp=%v hum=%v",
			m.Temperature, m.Humidity,
		)
	}
	if calls != 2 {
		t.Fatalf("expected failed metrics not to be cached, got %d calls", calls)
	}
}
