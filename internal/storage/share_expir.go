// Copyright (C) 2021-2022 Leonid Maslakov.

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

package storage


import(
	"errors"
)

var(
	ErrUnknowExpir = errors.New("unknown expiration")
)

func ExpirationToTime(expiration string) (int64, error) {
	switch expiration {
	// Minutes
	case "5m":
		return 60 * 5, nil

	case "10m":
		return 60 * 10, nil

	case "20m":
		return 60 * 20, nil

	case "30m":
		return 60 * 30, nil

	case "40m":
		return 60 * 40, nil

	case "50m":
		return 60 * 50, nil

	// Hours
	case "1h":
		return 60 * 60, nil

	case "2h":
		return 60 * 60 * 2, nil

	case "4h":
		return 60 * 60 * 4, nil

	case "12h":
		return 60 * 60 * 12, nil

	// Days
	case "1d":
		return 60 * 60 * 24, nil

	case "2d":
		return 60 * 60 * 24 * 2, nil

	case "3d":
		return 60 * 60 * 24 * 3, nil

	case "4d":
		return 60 * 60 * 24 * 4, nil

	case "5d":
		return 60 * 60 * 24 * 5, nil

	case "6d":
		return 60 * 60 * 24 * 6, nil

	// Weeks
	case "1w":
		return 60 * 60 * 24 * 7, nil

	case "2w":
		return 60 * 60 * 24 * 7 * 2, nil

	case "3w":
		return 60 * 60 * 24 * 7 * 3, nil
	}

	return 0, ErrUnknowExpir
}
