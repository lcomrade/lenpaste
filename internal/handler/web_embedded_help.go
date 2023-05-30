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
	"net"
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

// Pattern: /emb_help/
func (hand *handler) embeddedHelpHand(c *gin.Context) {
	type embHelpTmpl struct {
		ID         string
		DeleteTime int64
		OneUse     bool

		Protocol string
		Host     string

		Translate func(string, ...interface{}) template.HTML
		Highlight func(string, string) template.HTML
	}

	// Check rate limit
	err := hand.db.RateLimitCheck(model.RLPasteGet, net.IP(c.ClientIP()))
	if err != nil {
		hand.writeErrorWeb(c, err)
		return
	}

	// Get paste ID
	pasteID := c.Param("id")

	// Read DB
	paste, err := hand.db.PasteGet(pasteID)
	if err != nil {
		hand.writeErrorWeb(c, err)
		return
	}

	// Show paste
	tmplData := embHelpTmpl{
		ID:         paste.ID,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,
		Protocol:   getProtocol(c),
		Host:       getHost(c),
		Translate:  hand.l10n.findLocale(c).translate,
		Highlight:  hand.themes.findTheme(c, hand.cfg.UI.DefaultTheme).tryHighlight,
	}

	c.HTML(http.StatusOK, "emb_help.tmpl", tmplData)
}
