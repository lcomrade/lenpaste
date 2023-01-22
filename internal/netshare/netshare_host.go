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
	"net"
	"net/http"
	"strings"
)

func GetHost(req *http.Request) string {
	// Read header
	xHost := req.Header.Get("X-Forwarded-Host")

	// Check
	if xHost != "" {
		return xHost
	}

	return req.Host
}

func GetProtocol(req *http.Request) string {
	// X-Forwarded-Proto
	xProto := req.Header.Get("X-Forwarded-Proto")

	if xProto != "" {
		return xProto
	}

	// Else real protocol
	return req.URL.Scheme
}

func GetClientAddr(req *http.Request) net.IP {
	// X-Real-IP
	xReal := req.Header.Get("X-Real-IP")
	if xReal != "" {
		return net.ParseIP(xReal)
	}

	// X-Forwarded-For
	xFor := req.Header.Get("X-Forwarded-For")
	xFor = strings.Split(xFor, ",")[0]

	if xFor != "" {
		return net.ParseIP(xFor)
	}

	// Else use real client address
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil
	}

	return net.ParseIP(host)
}
