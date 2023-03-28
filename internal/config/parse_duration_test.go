package config

import (
	"testing"
)

func TestParseDuration(t *testing.T) {
	testData := map[string]int64{
		"10m":   60 * 10,
		"1h 1d": 60 * 60 * 25,
		"1w":    60 * 60 * 24 * 7,
		"365d":  60 * 60 * 24 * 365,
	}

	for s, exp := range testData {
		res, err := parseDuration(s)
		if err != nil {
			t.Fatal(err)
		}

		if exp != res {
			t.Error(exp, "!=", res)
		}
	}
}
