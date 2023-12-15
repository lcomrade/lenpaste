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
	"errors"
	"strconv"
	"time"
)

func parseDuration(s string) (time.Duration, error) {
	var out int64

	var tmp string
	for _, c := range s {
		if c == ' ' {
			continue
		}

		if '0' <= c && c <= '9' {
			tmp += string(c)
			continue
		}

		val, err := strconv.ParseInt(tmp, 10, 64)
		if err != nil {
			return 0, errors.New("invalid format \"" + s + "\"")
		}

		switch c {
		case 'm':
			out += val * 60
		case 'h':
			out += val * 60 * 60
		case 'd':
			out += val * 60 * 60 * 24
		case 'w':
			out += val * 60 * 60 * 24 * 7
		default:
			return 0, errors.New("invalid format \"" + s + "\"")
		}

		tmp = ""
	}

	return time.Duration(out) * time.Second, nil
}
