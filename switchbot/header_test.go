package switchbot_test

import (
	"roomctl/switchbot"
	"testing"
)

func TestHmacSHA256(t *testing.T) {
	s := switchbot.HmacSHA256("secret key", "This is a pen.")

	if s != "At/0AFUMncSeNiymj4gitbR5bhPNLe1lajj+gVIL+/c=" {
		t.Error("fail")
	}
}
