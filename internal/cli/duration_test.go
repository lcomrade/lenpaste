// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.

package cli

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	testData := map[string]time.Duration{
		"10m":   60 * 10 * time.Second,
		"1h 1d": 60 * 60 * 25 * time.Second,
		"1h1d": 60 * 60 * 25 * time.Second,
		"1w":    60 * 60 * 24 * 7 * time.Second,
		"365d":  60 * 60 * 24 * 365 * time.Second,
	}

	for s, exp := range testData {
		res, err := parseDuration(s)
		if err != nil {
			t.Fatal(err)
		}

		if exp != res {
			t.Error("expected", exp, "but got", res, "(input:", s, ")")
		}
	}
}
