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

package raw

import (
	"io"
	"net/http"
)

func (data Data) errorNotFound(rw http.ResponseWriter, req *http.Request) {
	// Write response header
	rw.Header().Set("Content-type", "text/plain")
	rw.WriteHeader(404)

	// Send
	_, e := io.WriteString(rw, "404 Not Found")
	if e != nil {
		data.Log.HttpError(req, e)
	}
}

func (data Data) errorInternal(rw http.ResponseWriter, req *http.Request, err error) {
	// Write to log
	data.Log.HttpError(req, err)

	// Write response header
	rw.Header().Set("Content-type", "text/plain")
	rw.WriteHeader(500)

	// Send
	_, e := io.WriteString(rw, "500 Internal Server Error")
	if e != nil {
		data.Log.HttpError(req, e)
	}
}
