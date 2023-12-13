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

package apiv1

import (
	"encoding/json"
	"errors"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"github.com/lcomrade/lenpaste/internal/storage"
	"net/http"
	"strconv"
)

type errorType struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	var resp errorType

	var eTmp429 *netshare.ErrTooManyRequests

	if e == netshare.ErrBadRequest {
		resp.Code = 400
		resp.Error = "Bad Request"

	} else if e == netshare.ErrUnauthorized {
		rw.Header().Add("WWW-Authenticate", "Basic")
		resp.Code = 401
		resp.Error = "Unauthorized"

	} else if e == storage.ErrNotFoundID {
		resp.Code = 404
		resp.Error = "Could not find ID"

	} else if e == netshare.ErrNotFound {
		resp.Code = 404
		resp.Error = "Not Found"

	} else if e == netshare.ErrMethodNotAllowed {
		resp.Code = 405
		resp.Error = "Method Not Allowed"

	} else if e == netshare.ErrPayloadTooLarge {
		resp.Code = 413
		resp.Error = "Payload Too Large"

	} else if errors.As(e, &eTmp429) {
		resp.Code = 429
		resp.Error = "Too Many Requests"
		rw.Header().Set("Retry-After", strconv.FormatInt(eTmp429.RetryAfter, 10))

	} else {
		resp.Code = 500
		resp.Error = "Internal Server Error"
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(resp.Code)

	err := json.NewEncoder(rw).Encode(resp)
	if err != nil {
		return 500, err
	}

	return resp.Code, nil
}
