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

package config

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	var out int64

	for _, part := range strings.Split(s, " ") {
		if strings.HasSuffix(part, "m") {
			val, err := strconv.Atoi(part[:len(part)-1])
			if err != nil {
				return 0, errors.New(`parse duration: invalid format "` + part + `"`)
			}
			out = out + (int64(val) * 60)
			continue
		}

		if strings.HasSuffix(part, "h") {
			val, err := strconv.Atoi(part[:len(part)-1])
			if err != nil {
				return 0, errors.New(`parse duration: invalid format "` + part + `"`)
			}
			out = out + (int64(val) * 60 * 60)
			continue
		}

		if strings.HasSuffix(part, "d") {
			val, err := strconv.Atoi(part[:len(part)-1])
			if err != nil {
				return 0, errors.New(`parse duration: invalid format "` + part + `"`)
			}
			out = out + int64((val)*60*60*24)
			continue
		}

		if strings.HasSuffix(part, "w") {
			val, err := strconv.Atoi(part[:len(part)-1])
			if err != nil {
				return 0, errors.New(`parse duration: invalid format "` + part + `"`)
			}
			out = out + int64((val)*60*60*24*7)
			continue
		}

		return 0, errors.New(`parse duration: invalid format "` + part + `"`)
	}

	return time.Duration(out) * time.Second, nil
}
