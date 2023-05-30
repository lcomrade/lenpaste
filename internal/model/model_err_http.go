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

package model

import (
	"strconv"
)

var (
	ErrBadRequest       = NewError(400, "Bad Request")
	ErrUnauthorized     = NewError(401, "Unauthorized")
	ErrNotFound         = NewError(404, "Not Found")
	ErrMethodNotAllowed = NewError(405, "Method Not Allowed")
	ErrPayloadTooLarge  = NewError(413, "Payload Too Large")
)

func ErrTooManyRequestsNew(retryAfter int64) error {
	header := make(map[string]string)
	header["Retry-After"] = strconv.FormatInt(retryAfter, 10)

	return NewErrorWithHeader(429, "Too Many Requests", header)
}
