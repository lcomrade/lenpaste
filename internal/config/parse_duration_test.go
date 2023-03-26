package config

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	testData := map[string]time.Duration{
		"10m":   time.Second * 60 * 10,
		"1h 1d": time.Second * 60 * 60 * 25,
		"1w":    time.Second * 60 * 60 * 24 * 7,
		"365d":  time.Second * 60 * 60 * 24 * 365,
	}

	for s, exp := range testData {
		res, err := ParseDuration(s)
		if err != nil {
			t.Fatal(err)
		}

		if exp != res {
			t.Error(exp, "!=", res)
		}
	}
}
