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
func (hand *handler) docsHand(c *gin.Context) {
	c.HTML(http.StatusOK, "docs.tmpl", docsTmpl{Translate: hand.l10n.findLocale(c).translate})
}

// Pattern: /docs/apiv1
func (hand *handler) docsApiV1Hand(c *gin.Context) {
	c.HTML(http.StatusOK, "docs_apiv1.tmpl", docsApiV1Tmpl{
		MaxLenAuthorAll: model.MaxLengthAuthorAll,
		Translate:       hand.l10n.findLocale(c).translate,
		Highlight:       hand.themes.findTheme(c, hand.cfg.UI.DefaultTheme).tryHighlight,
	})
}

// Pattern: /docs/api_libs
func (hand *handler) docsApiLibsHand(c *gin.Context) {
	c.HTML(http.StatusOK, "docs_api_libs.tmpl", docsTmpl{Translate: hand.l10n.findLocale(c).translate})
}
