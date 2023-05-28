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

package handler

import (
	"github.com/gin-gonic/gin"
)

const cookieMaxAge = 60 * 60 * 24 * 360 * 50 // 50 year

func setCookie(c *gin.Context, name, value string) {
	if value != "" {
		// Set cookie
		c.SetCookie(
			name, value,
			cookieMaxAge, "/", "",
			false, false,
		)

	} else {
		// Delete cookie
		c.SetCookie(
			name, "",
			-1, "/", "",
			false, false,
		)
	}
}

func getCookie(c *gin.Context, name string) string {
	val, err := c.Cookie(name)
	if err != nil {
		return ""
	}

	return val
}

func getHost(c *gin.Context) string {
	// Read header
	xHost := c.Request.Header.Get("X-Forwarded-Host")

	// Check
	if xHost != "" {
		return xHost
	}

	return c.Request.Host
}

func getProtocol(c *gin.Context) string {
	// X-Forwarded-Proto
	xProto := c.Request.Header.Get("X-Forwarded-Proto")

	if xProto != "" {
		return xProto
	}

	// Else real protocol
	return c.Request.URL.Scheme
}
