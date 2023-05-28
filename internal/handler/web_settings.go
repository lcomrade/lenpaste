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
	"html/template"
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

// Pattern: /settings
func (hand *handler) settingsHand(c *gin.Context) {
	type settingsTmpl struct {
		Language         string
		LanguageSelector map[string]string

		Theme         string
		ThemeSelector map[string]string

		AuthorAllMaxLen int
		Author          string
		AuthorEmail     string
		AuthorURL       string

		AuthOk bool

		Translate func(string, ...interface{}) template.HTML
	}

	if c.Request.Method != "POST" {
		// Prepare data
		dataTmpl := settingsTmpl{
			Language:         getCookie(c, "lang"),
			LanguageSelector: hand.l10n.names,
			Theme:            getCookie(c, "theme"),
			ThemeSelector:    hand.themes.getForLocale(c, hand.l10n),
			AuthorAllMaxLen:  model.MaxLengthAuthorAll,
			Author:           getCookie(c, "author"),
			AuthorEmail:      getCookie(c, "authorEmail"),
			AuthorURL:        getCookie(c, "authorURL"),
			Translate:        hand.l10n.findLocale(c).translate,
		}

		if dataTmpl.Theme == "" {
			dataTmpl.Theme = hand.cfg.UI.DefaultTheme
		}

		// Show page
		c.HTML(http.StatusOK, "settings.tmpl", dataTmpl)

		// Else update settings
	} else {

		setCookie(c, "lang", c.PostForm("lang"))
		setCookie(c, "theme", c.PostForm("theme"))

		setCookie(c, "author", c.PostForm("author"))
		setCookie(c, "authorEmail", c.PostForm("authorEmail"))
		setCookie(c, "authorURL", c.PostForm("authorURL"))

		c.Redirect(http.StatusFound, "/settings")
	}
}
