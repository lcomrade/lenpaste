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

package web

import "github.com/gin-gonic/gin"

func writeRedirect(c *gin.Context, newURL string, code int) {
	if newURL == "" {
		newURL = "/"
	}

	if req.URL.RawQuery != "" {
		newURL = newURL + "?" + req.URL.RawQuery
	}

	c.Header("Location", newURL)
	rw.WriteHeader(code)
}
