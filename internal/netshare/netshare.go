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

package netshare

import (
	"errors"
)

const (
	MaxLengthAuthorAll = 100 // Max length or paste author name, email and URL.
)

var (
	ErrBadRequest       = errors.New("Bad Request")        // 400
	ErrUnauthorized     = errors.New("Unauthorized")       // 401
	ErrNotFound         = errors.New("Not Found")          // 404
	ErrMethodNotAllowed = errors.New("Method Not Allowed") // 405
	ErrPayloadTooLarge  = errors.New("Payload Too Large")  // 413
	//	ErrTooManyRequests  = errors.New("Too Many Requests")     // 429
	ErrInternal = errors.New("Internal Server Error") // 500
)

type ErrTooManyRequests struct {
	s          string
	RetryAfter int64
}

func (e *ErrTooManyRequests) Error() string {
	return e.s
}

func ErrTooManyRequestsNew(retryAfter int64) *ErrTooManyRequests {
	return &ErrTooManyRequests{
		s:          "Too Many Requests",
		RetryAfter: retryAfter,
	}
}
