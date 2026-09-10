package switchbot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Body       body   `json:"body"`
}

const switchBotAPISuccess = 100

var (
	switchBotHTTPRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "roomctl",
		Subsystem: "switchbot",
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests to the SwitchBot API.",
	}, []string{"code", "method"})
	switchBotHTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "roomctl",
		Subsystem: "switchbot",
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests to the SwitchBot API.",
	}, []string{"code", "method"})
	switchBotHTTPInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "roomctl",
		Subsystem: "switchbot",
		Name:      "http_in_flight_requests",
		Help:      "Current number of in-flight HTTP requests to the SwitchBot API.",
	})
	switchBotHTTPClient = &http.Client{
		Transport: promhttp.InstrumentRoundTripperInFlight(
			switchBotHTTPInFlight,
			promhttp.InstrumentRoundTripperDuration(
				switchBotHTTPRequestDuration,
				promhttp.InstrumentRoundTripperCounter(
					switchBotHTTPRequestTotal,
					http.DefaultTransport,
				),
			),
		),
	}
)

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

	resp, err := switchBotHTTPClient.Do(r)
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
