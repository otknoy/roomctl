package mflight_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"roomctl/mflight"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSensorMonitor(t *testing.T) {
	s := newStubServer(t, `
      <db>
        <table id="67243">
          <time>202101030000</time>
          <unixtime>1609599600</unixtime>
          <temp>22.0</temp>
          <humi>43.3</humi>
          <illu>405</illu>
        </table>
        <table id="67244">
          <time>202101030005</time>
          <unixtime>1609599900</unixtime>
          <temp>21.9</temp>
          <humi>43.0</humi>
          <illu>406</illu>
        </table>
      </db>`)
	defer s.Close()

	c := &mflight.ClientImpl{
		BaseUrl:  s.URL,
		MobileId: "test-mobile-id",
	}

	temp, humi, illu, err := c.GetMetrics(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(float32(22.0), temp); diff != "" {
		t.Errorf("response differs.\n%v", diff)
	}
	if diff := cmp.Diff(float32(43.3), humi); diff != "" {
		t.Errorf("response differs.\n%v", diff)
	}
	if diff := cmp.Diff(uint(405), illu); diff != "" {
		t.Errorf("response differs.\n%v", diff)
	}
}

func newStubServer(t *testing.T, response string) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path := r.URL.Path; path != "/SensorMonitorV2.xml" {
			t.Fatal(path)
		}

		if qs := r.URL.RawQuery; qs != "x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE=" {
			t.Fatal(qs)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	}))

	return s
}
