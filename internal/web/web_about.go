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

type aboutTmpl struct {
	Version     string
	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsExist bool

	AdminName string
	AdminMail string

	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

type aboutMinTmp struct {
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /about
func (data *Data) AboutHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Prepare data
	dataTmpl := aboutTmpl{
		Version:          *data.Version,
		TitleMaxLen:      *data.TitleMaxLen,
		BodyMaxLen:       *data.BodyMaxLen,
		MaxLifeTime:      *data.MaxLifeTime,
		ServerAbout:      *data.ServerAbout,
		ServerRules:      *data.ServerRules,
		ServerTermsExist: *data.ServerTermsExist,
		AdminName:        *data.AdminName,
		AdminMail:        *data.AdminMail,
		Highlight:        data.Themes.findTheme(req).tryHighlight,
		Translate:        data.Locales.findLocale(req).translate,
	}

	// Show page
	rw.Header().Set("Content-Type", "text/html")

	err := data.About.Execute(rw, dataTmpl)
	if err != nil {
		data.writeError(rw, req, err)
	}
}

// Pattern: /about/authors
func (data *Data) AuthorsHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Show page
	rw.Header().Set("Content-Type", "text/html")

	err := data.Authors.Execute(rw, aboutMinTmp{Translate: data.Locales.findLocale(req).translate})
	if err != nil {
		data.writeError(rw, req, err)
	}
}

// Pattern: /about/license
func (data *Data) LicenseHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Show page
	rw.Header().Set("Content-Type", "text/html")

	err := data.License.Execute(rw, aboutMinTmp{Translate: data.Locales.findLocale(req).translate})
	if err != nil {
		data.writeError(rw, req, err)
	}
}

// Pattern: /about/source_code
func (data *Data) SourceCodePageHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Show page
	rw.Header().Set("Content-Type", "text/html")

	err := data.SourceCodePage.Execute(rw, aboutMinTmp{Translate: data.Locales.findLocale(req).translate})
	if err != nil {
		data.writeError(rw, req, err)
	}
}
