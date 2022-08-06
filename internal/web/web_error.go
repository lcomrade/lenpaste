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
	"html/template"
	"net/http"
)

type errorTmpl struct {
	Code      int
	Error     string
	AdminName string
	AdminMail string
	Translate func(string, ...interface{}) template.HTML
}

func (data Data) errorBadRequest(rw http.ResponseWriter, req *http.Request) {
	// Write response header
	rw.Header().Set("Content-type", "text/html")
	rw.WriteHeader(400)

	// Render template
	errData := errorTmpl{
		Error:     "Bad Request",
		Code:      400,
		AdminName: *data.AdminName,
		AdminMail: *data.AdminMail,
		Translate: data.Locales.findLocale(req).translate,
	}

	e := data.ErrorPage.Execute(rw, errData)
	if e != nil {
		data.Log.HttpError(req, e)
	}
}

func (data Data) errorNotFound(rw http.ResponseWriter, req *http.Request) {
	// Write response header
	rw.Header().Set("Content-type", "text/html")
	rw.WriteHeader(404)

	// Render template
	errData := errorTmpl{
		Error:     "Not Found",
		Code:      404,
		AdminName: *data.AdminName,
		AdminMail: *data.AdminMail,
		Translate: data.Locales.findLocale(req).translate,
	}

	e := data.ErrorPage.Execute(rw, errData)
	if e != nil {
		data.Log.HttpError(req, e)
	}
}

func (data Data) errorInternal(rw http.ResponseWriter, req *http.Request, err error) {
	// Write to log
	data.Log.HttpError(req, err)

	// Write response header
	rw.Header().Set("Content-type", "text/html")
	rw.WriteHeader(500)

	// Render template
	errData := errorTmpl{
		Error:     "Internal Server Error",
		Code:      500,
		AdminName: *data.AdminName,
		AdminMail: *data.AdminMail,
		Translate: data.Locales.findLocale(req).translate,
	}

	e := data.ErrorPage.Execute(rw, errData)
	if e != nil {
		data.Log.HttpError(req, e)
	}
}
