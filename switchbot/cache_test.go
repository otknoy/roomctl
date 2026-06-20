package switchbot

import (
	"context"
	"errors"
	"testing"
)

type clientFunc func(context.Context) (float32, float32, error)

func (f clientFunc) GetMetrics(ctx context.Context) (float32, float32, error) {
	return f(ctx)
}

func TestCacheStoresSuccessfulMetrics(t *testing.T) {
	calls := 0
	client := clientFunc(func(context.Context) (float32, float32, error) {
		calls++
		if calls == 1 {
			return 23.5, 48, nil
		}

		return 99, 99, nil
	})
	cache := NewCacheClient(client)

	temp, hum, err := cache.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected fresh metrics, got %v", err)
	}
	if temp != 23.5 || hum != 48 {
		t.Fatalf("expected fresh metrics, got temp=%v hum=%v", temp, hum)
	}

	temp, hum, err = cache.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected cached metrics, got %v", err)
	}
	if temp != 23.5 || hum != 48 {
		t.Fatalf("expected cached metrics, got temp=%v hum=%v", temp, hum)
	}
	if calls != 1 {
		t.Fatalf("expected successful metrics to be cached, got %d calls", calls)
	}
}

func TestCacheDoesNotStoreFailedMetrics(t *testing.T) {
	temporaryErr := errors.New("temporary error")
	calls := 0
	client := clientFunc(func(context.Context) (float32, float32, error) {
		calls++
		if calls == 1 {
			return 0, 0, temporaryErr
		}

		return 24.5, 51, nil
	})
	cache := NewCacheClient(client)

	temp, hum, err := cache.GetMetrics(context.Background())
	if !errors.Is(err, temporaryErr) {
		t.Fatalf("expected temporary error, got %v", err)
	}
	if temp != 0 || hum != 0 {
		t.Fatalf("failed response should return zero values, got temp=%v hum=%v", temp, hum)
	}

	temp, hum, err = cache.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	if temp != 24.5 || hum != 51 {
		t.Fatalf("expected fresh metrics after failure, got temp=%v hum=%v", temp, hum)
	}
	if calls != 2 {
		t.Fatalf("expected failed metrics not to be cached, got %d calls", calls)
	}
}
