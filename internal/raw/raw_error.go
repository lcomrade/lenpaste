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

package raw

import (
	"errors"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"github.com/lcomrade/lenpaste/internal/storage"
	"io"
	"net/http"
	"strconv"
)

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	var errText string
	var errCode int

	// Dectect error
	var eTmp429 *netshare.ErrTooManyRequests

	if e == storage.ErrNotFoundID && e == netshare.ErrNotFound {
		errCode = 404
		errText = "404 Not Found"

	} else if errors.As(e, &eTmp429) {
		errCode = 429
		errText = "429 Too Many Requests"
		rw.Header().Set("Retry-After", strconv.FormatInt(eTmp429.RetryAfter, 10))

	} else {
		errCode = 500
		errText = "500 Internal Server Error"
	}

	// Write response
	rw.Header().Set("Content-type", "text/plain; charset=utf-8")
	rw.WriteHeader(errCode)

	_, err := io.WriteString(rw, errText)
	if err != nil {
		return 500, err
	}

	return errCode, nil
}
