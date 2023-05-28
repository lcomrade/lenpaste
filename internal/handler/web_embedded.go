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
	"time"

	"git.lcomrade.su/root/lenpaste/internal/storage"
	"github.com/gin-gonic/gin"
)

type embTmpl struct {
	ID            string
	CreateTimeStr string
	DeleteTime    int64
	OneUse        bool
	Title         string
	Body          template.HTML

	ErrorNotFound bool
	Translate     func(string, ...interface{}) template.HTML
}

// Pattern: /emb/
func (hand *handler) embeddedHand(c *gin.Context) {
	errorNotFound := false

	// Check rate limit
	err := hand.db.RateLimitCheck("paste_get", net.IP(c.ClientIP()))
	if err != nil {
		hand.writeErrorWeb(c, err)
		return
	}

	// Get paste ID
	pasteID := c.Param("id")

	// Read DB
	paste, err := hand.db.PasteGet(pasteID)
	if err != nil {
		if err == storage.ErrNotFoundID {
			errorNotFound = true

		} else {
			hand.writeErrorWeb(c, err)
			return
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()

	tmplData := embTmpl{
		ID:            paste.ID,
		CreateTimeStr: createTime.Format("1 Jan, 2006"),
		DeleteTime:    paste.DeleteTime,
		OneUse:        paste.OneUse,
		Title:         paste.Title,
		Body:          tryHighlight(paste.Body, paste.Syntax, "monokai"),

		ErrorNotFound: errorNotFound,
		Translate:     hand.l10n.findLocale(c).translate,
	}

	// Show paste
	c.HTML(http.StatusOK, "emb.tmpl", tmplData)
}
