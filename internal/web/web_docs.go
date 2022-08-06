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

type docsTmpl struct {
	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /docs
func (data Data) DocsHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "text/html")
	err := data.Docs.Execute(rw, docsTmpl{Translate: data.Locales.findLocale(req).translate})
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}

// Pattern: /docs/apiv1
func (data Data) DocsApiV1Hand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "text/html")
	err := data.DocsApiV1.Execute(rw, docsTmpl{
		Translate: data.Locales.findLocale(req).translate,
		Highlight: tryHighlight,
	})
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}

// Pattern: /docs/api_libs
func (data Data) DocsApiLibsHand(rw http.ResponseWriter, req *http.Request) {
	data.Log.HttpRequest(req)

	rw.Header().Set("Content-Type", "text/html")
	err := data.DocsApiLibs.Execute(rw, docsTmpl{Translate: data.Locales.findLocale(req).translate})
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
