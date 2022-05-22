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

package web

import(
	"net/http"
)

// Pattern: /docs
func (data Data) DocsHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "text/html")
	data.Docs.Execute(rw, "")
}

// Pattern: /docs/apiv1
func (data Data) DocsAPIV1Hand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "text/html")
	data.DocsApiV1.Execute(rw, "")
}
