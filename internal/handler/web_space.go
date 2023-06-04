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

	"github.com/gin-gonic/gin"
)

func (hand *handler) spaceCSSHand(c *gin.Context) {
	c.Header("Content-Type", "text/css; charset=utf-8")
	err := hand.spaceCSS.Execute(c.Writer, jsTmpl{
		Translate: hand.l10n.findLocale(c).translate,
		Theme:     hand.themes.findTheme(c, hand.cfg.UI.DefaultTheme).theme,
	})

	if err != nil {
		hand.writeErrorWeb(c, err)
		return
	}
}

func (hand *handler) spaceJSHand(c *gin.Context) {
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	err := hand.spaceJS.Execute(c.Writer, jsTmpl{
		Translate: hand.l10n.findLocale(c).translate,
		Theme:     hand.themes.findTheme(c, hand.cfg.UI.DefaultTheme).theme,
	})

	if err != nil {
		hand.writeErrorWeb(c, err)
		return
	}
}

func (hand *handler) spaceHand(c *gin.Context) {
	type spaceTmpl struct {
		Translate func(string, ...interface{}) template.HTML
		Theme     func(string) string
	}

	c.HTML(http.StatusOK, "space.tmpl", spaceTmpl{
		Translate: hand.l10n.findLocale(c).translate,
		Theme:     hand.themes.findTheme(c, hand.cfg.UI.DefaultTheme).theme,
	})
}