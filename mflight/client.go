package mflight

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	GetMetrics(ctx context.Context) (temp, hum float32, illu uint, err error)
}

var _ Client = (*ClientImpl)(nil)

type ClientImpl struct {
	BaseUrl  string
	MobileId string
}

func (c *ClientImpl) GetMetrics(ctx context.Context) (temp, hum float32, illu uint, err error) {
	r := buildRequestWithContext(ctx, c.BaseUrl, c.MobileId)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()

	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, err
	}

	var res response
	if err := xml.Unmarshal(byteArray, &res); err != nil {
		return 0, 0, 0, err
	}

	tables := res.Tables
	if len(tables) == 0 {
		return 0, 0, 0, errors.New("empty response")
	}

	t := tables[0]

	return t.Temperature, t.Humidity, uint(t.Illuminance), nil
}

func buildRequestWithContext(ctx context.Context, baseURL, mobileID string) *http.Request {
	qs := url.Values{
		"x-KEY_MOBILE_ID":   []string{mobileID},
		"x-KEY_UPDATE_DATE": []string{""},
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	req.URL.Path = "/SensorMonitorV2.xml"
	req.URL.RawQuery = qs.Encode()

	return req
}

type response struct {
	Tables []Table `xml:"table"`
}

type Table struct {
	ID          int64   `xml:"id,attr"`
	Time        string  `xml:"time"`
	Unixtime    int64   `xml:"unixtime"`
	Temperature float32 `xml:"temp"`
	Humidity    float32 `xml:"humi"`
	Illuminance int16   `xml:"illu"`
}
