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

	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"github.com/gin-gonic/gin"
)

type embHelpTmpl struct {
	ID         string
	DeleteTime int64
	OneUse     bool

	Protocol string
	Host     string

	Translate func(string, ...interface{}) template.HTML
	Highlight func(string, string) template.HTML
}

// Pattern: /emb_help/
func (hand *handler) embeddedHelpHand(c *gin.Context) {
	// Check rate limit
	err := hand.db.RateLimitCheck("paste_get", netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[10:])

	// Read DB
	paste, err := hand.db.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// Show paste
	tmplData := embHelpTmpl{
		ID:         paste.ID,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,
		Protocol:   netshare.GetProtocol(req),
		Host:       netshare.GetHost(req),
		Translate:  hand.l10n.findLocale(req).translate,
		Highlight:  data.themes.findTheme(req, hand.cfg.UI.DefaultTheme).tryHighlight,
	}

	return data.embeddedHelpPage.Execute(rw, tmplData)
}
