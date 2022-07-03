package switchbot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type response struct {
	Body body
}

type body struct {
	Humidity    float32
	Temparature float32
}

func GetMetrics(ctx context.Context, deviceId string, token string) (hum, temp float32, err error) {
	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.switch-bot.com/v1.0/devices/%s/status", deviceId),
		nil,
	)
	if err != nil {
		return 0.0, 0.0, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0.0, 0.0, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, 0.0, err
	}

	var res response
	if err = json.Unmarshal(bytes, &res); err != nil {
		return 0.0, 0.0, err
	}

	return res.Body.Humidity, res.Body.Temparature, nil
}
