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
func (data *Data) aboutHand(rw http.ResponseWriter, req *http.Request) error {
	lang := data.l10n.detectLanguage(req)

	dataTmpl := aboutTmpl{
		Version:          model.Version,
		TitleMaxLen:      data.cfg.Paste.TitleMaxLen,
		BodyMaxLen:       data.cfg.Paste.BodyMaxLen,
		MaxLifeTime:      data.cfg.Paste.MaxLifetime,
		ServerAbout:      data.cfg.GetAbout(lang),
		ServerRules:      data.cfg.GetRules(lang),
		ServerTermsExist: data.cfg.TermsOfUse != nil,
		AdminName:        data.cfg.Public.AdminName,
		AdminMail:        data.cfg.Public.AdminMail,
		Highlight:        data.themes.findTheme(req, data.cfg.UI.DefaultTheme).tryHighlight,
		Translate:        data.l10n.findLocale(req).translate,
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.about.Execute(rw, dataTmpl)
}

// Pattern: /about/authors
func (data *Data) authorsHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.authors.Execute(rw, aboutMinTmp{Translate: data.l10n.findLocale(req).translate})
}

// Pattern: /about/license
func (data *Data) licenseHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.license.Execute(rw, aboutMinTmp{Translate: data.l10n.findLocale(req).translate})
}

// Pattern: /about/source_code
func (data *Data) sourceCodePageHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.sourceCodePage.Execute(rw, aboutMinTmp{Translate: data.l10n.findLocale(req).translate})
}
