package switchbot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Body       body   `json:"body"`
}

const switchBotAPISuccess = 100

type body struct {
	Temperature float32
	Humidity    float32
}

type Client interface {
	GetMetrics(ctx context.Context) (*Metrics, error)
}

var _ Client = (*ClientImpl)(nil)

type ClientImpl struct {
	Token    string
	Secret   string
	DeviceId string
}

type Metrics struct {
	Temperature float32
	Humidity    float32
}

func (c *ClientImpl) GetMetrics(ctx context.Context) (*Metrics, error) {
	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.switch-bot.com/v1.1/devices/%s/status", c.DeviceId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	r.Header = makeHeader(c.Token, c.Secret)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	return decodeMetrics(resp.Body)
}

func decodeMetrics(r io.Reader) (*Metrics, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var res response
	if err = json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}
	if res.StatusCode != switchBotAPISuccess {
		return nil, fmt.Errorf("switchbot API status %d: %s", res.StatusCode, res.Message)
	}

	return &Metrics{
		Temperature: res.Body.Temperature,
		Humidity:    res.Body.Humidity,
	}, nil
}
