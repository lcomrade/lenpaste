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

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
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
func (hand *handler) aboutHand(c *gin.Context) {
	lang := hand.l10n.detectLanguage(req)

	dataTmpl := aboutTmpl{
		Version:          model.Version,
		TitleMaxLen:      hand.cfg.Paste.TitleMaxLen,
		BodyMaxLen:       hand.cfg.Paste.BodyMaxLen,
		MaxLifeTime:      hand.cfg.Paste.MaxLifetime,
		ServerAbout:      hand.cfg.GetAbout(lang),
		ServerRules:      hand.cfg.GetRules(lang),
		ServerTermsExist: hand.cfg.TermsOfUse != nil,
		AdminName:        hand.cfg.Public.AdminName,
		AdminMail:        hand.cfg.Public.AdminMail,
		Highlight:        data.themes.findTheme(req, hand.cfg.UI.DefaultTheme).tryHighlight,
		Translate:        hand.l10n.findLocale(req).translate,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	return data.about.Execute(rw, dataTmpl)
}

// Pattern: /about/authors
func (hand *handler) authorsHand(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	return data.authors.Execute(rw, aboutMinTmp{Translate: hand.l10n.findLocale(req).translate})
}

// Pattern: /about/license
func (hand *handler) licenseHand(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	return data.license.Execute(rw, aboutMinTmp{Translate: hand.l10n.findLocale(req).translate})
}

// Pattern: /about/source_code
func (hand *handler) sourceCodePageHand(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	return data.sourceCodePage.Execute(rw, aboutMinTmp{Translate: hand.l10n.findLocale(req).translate})
}
