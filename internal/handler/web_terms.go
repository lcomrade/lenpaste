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

// Pattern: /terms
func (hand *handler) termsOfUseHand(c *gin.Context) {
	type termsOfUseTmpl struct {
		TermsOfUse string

		Highlight func(string, string) template.HTML
		Translate func(string, ...interface{}) template.HTML
	}

	c.HTML(http.StatusOK, "terms.tmpl", termsOfUseTmpl{
		TermsOfUse: hand.cfg.GetTermsOfUse(hand.l10n.detectLanguage(c)),
		Highlight:  hand.themes.findTheme(c, hand.cfg.UI.DefaultTheme).tryHighlight,
		Translate:  hand.l10n.findLocale(c).translate},
	)
}
