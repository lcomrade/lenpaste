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

	"github.com/gin-gonic/gin"
)

type termsOfUseTmpl struct {
	TermsOfUse string

	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /terms
func (hand *handler) termsOfUseHand(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	return data.termsOfUse.Execute(rw, termsOfUseTmpl{
		TermsOfUse: hand.cfg.GetTermsOfUse(hand.l10n.detectLanguage(req)),
		Highlight:  data.themes.findTheme(req, hand.cfg.UI.DefaultTheme).tryHighlight,
		Translate:  hand.l10n.findLocale(req).translate},
	)
}
