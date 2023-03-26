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
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	// Parse error
	resp := model.ParseError(e)

	if resp.Code > 499 {
		data.log.HttpError(req, e)
	}

	// Prepare header
	rw.Header().Set("Content-Type", "application/json")
	for key, val := range resp.Header {
		rw.Header().Set(key, val)
	}

	// Set HTTP status code
	rw.WriteHeader(resp.Code)

	// Send JSON response body
	err := json.NewEncoder(rw).Encode(resp)
	if err != nil {
		return 500, err
	}

	return resp.Code, nil
}
