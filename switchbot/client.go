package switchbot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type response struct {
	Body body
}

type body struct {
	Temperature float32
	Humidity    float32
}

func GetMetrics(ctx context.Context, deviceId string, token string) (temp, hum float32, err error) {
	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.switch-bot.com/v1.0/devices/%s/status", deviceId),
		nil,
	)
	if err != nil {
		return 0.0, 0.0, err
	}

	r.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0.0, 0.0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0.0, 0.0, fmt.Errorf("status: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, 0.0, err
	}

	log.Println(string(bytes))

	var res response
	if err = json.Unmarshal(bytes, &res); err != nil {
		return 0.0, 0.0, err
	}

	log.Println(res)

	return res.Body.Temperature, res.Body.Humidity, nil
}
