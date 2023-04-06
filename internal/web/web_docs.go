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

import (
	"html/template"
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

type docsTmpl struct {
	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

type docsApiV1Tmpl struct {
	MaxLenAuthorAll int

	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /docs
func (data *Data) docsHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.docs.Execute(rw, docsTmpl{Translate: data.l10n.findLocale(req).translate})
}

// Pattern: /docs/apiv1
func (data *Data) docsApiV1Hand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.docsApiV1.Execute(rw, docsApiV1Tmpl{
		MaxLenAuthorAll: model.MaxLengthAuthorAll,
		Translate:       data.l10n.findLocale(req).translate,
		Highlight:       data.themes.findTheme(req, data.cfg.UI.DefaultTheme).tryHighlight,
	})
}

// Pattern: /docs/api_libs
func (data *Data) docsApiLibsHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.docsApiLibs.Execute(rw, docsTmpl{Translate: data.l10n.findLocale(req).translate})
}
