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

package apiv1

import (
	"encoding/json"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"net/http"
)

type errorType struct {
	Code  int
	Error string
}

func (data Data) writeError(rw http.ResponseWriter, req *http.Request, err error) {
	var resp errorType

	if err == errBadRequest {
		resp.Code = 400
		resp.Error = "Bad Request"

	} else if err == storage.ErrNotFoundID {
		resp.Code = 403
		resp.Error = "Could not find ID"

	} else {
		resp.Code = 500
		resp.Error = "Internal Server Error"
		data.Log.HttpError(req, err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(resp.Code)

	e := json.NewEncoder(rw).Encode(resp)
	if err != nil {
		data.Log.HttpError(req, e)
		return
	}
}
