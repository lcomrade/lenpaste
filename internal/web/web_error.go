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

import (
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"html/template"
	"net/http"
	"strconv"
)

type errorTmpl struct {
	Code      int
	AdminName string
	AdminMail string
	Translate func(string, ...interface{}) template.HTML
}

func (data Data) writeError(rw http.ResponseWriter, req *http.Request, e error) {
	errData := errorTmpl{
		Code:      0,
		AdminName: *data.AdminName,
		AdminMail: *data.AdminMail,
		Translate: data.Locales.findLocale(req).translate,
	}

	// Dectect error
	switch e {
	case netshare.ErrBadRequest:
		errData.Code = 400
	case netshare.ErrUnauthorized:
		errData.Code = 401
	case storage.ErrNotFoundID:
		errData.Code = 404
	case netshare.ErrNotFound:
		errData.Code = 404
	case netshare.ErrMethodNotAllowed:
		errData.Code = 405
	case netshare.ErrPayloadTooLarge:
		errData.Code = 413
	case netshare.ErrTooManyRequests:
		errData.Code = 429
		rw.Header().Set("Retry-After", strconv.Itoa(netshare.RateLimitPeriod))
	default:
		errData.Code = 500
	}

	// Log Internal Server Error if need
	if errData.Code >= 500 && e != netshare.ErrInternal {
		data.Log.HttpError(req, e)
	}

	// Write response header
	rw.Header().Set("Content-type", "text/html")
	rw.WriteHeader(errData.Code)

	// Render template
	err := data.ErrorPage.Execute(rw, errData)
	if err != nil {
		data.Log.HttpError(req, err)
	}
}
