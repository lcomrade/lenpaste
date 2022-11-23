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

func GetProtocol(header http.Header) string {
	// Read header
	xProto := header.Get("X-Forwarded-Proto")

	// Check
	if xProto != "" {
		return xProto
	}

	return "http"
}

func GetClientAddr(req *http.Request) net.IP {
	// Read header
	xFor := req.Header.Get("X-Forwarded-For")
	xFor = strings.Split(xFor, ",")[0]

	// Check
	if xFor != "" {
		return net.ParseIP(xFor)
	}

	// Else use client address
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil
	}

	return net.ParseIP(host)
}
