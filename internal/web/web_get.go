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
	"time"

	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"git.lcomrade.su/root/lineend"
	"github.com/gin-gonic/gin"
)

type pasteTmpl struct {
	ID         string
	Title      string
	Body       template.HTML
	Syntax     string
	CreateTime int64
	DeleteTime int64
	OneUse     bool

	LineEnd       string
	CreateTimeStr string
	DeleteTimeStr string

	Author      string
	AuthorEmail string
	AuthorURL   string

	Translate func(string, ...interface{}) template.HTML
}

type pasteContinueTmpl struct {
	ID        string
	Translate func(string, ...interface{}) template.HTML
}

func (hand *handler) getPasteHand(c *gin.Context) {
	// Check rate limit
	err := hand.db.RateLimitCheck("paste_get", netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[1:])

	// Read DB
	paste, err := hand.db.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse {
		// If continue button not pressed
		req.ParseForm()

		if req.PostForm.Get("oneUseContinue") != "true" {
			tmplData := pasteContinueTmpl{
				ID:        paste.ID,
				Translate: hand.l10n.findLocale(req).translate,
			}

			return data.pasteContinue.Execute(rw, tmplData)
		}

		// If continue button pressed delete paste
		err = hand.db.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()
	deleteTime := time.Unix(paste.DeleteTime, 0).UTC()

	tmplData := pasteTmpl{
		ID:         paste.ID,
		Title:      paste.Title,
		Body:       data.themes.findTheme(req, hand.cfg.UI.DefaultTheme).tryHighlight(paste.Body, paste.Syntax),
		Syntax:     paste.Syntax,
		CreateTime: paste.CreateTime,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,

		CreateTimeStr: createTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		DeleteTimeStr: deleteTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),

		Author:      paste.Author,
		AuthorEmail: paste.AuthorEmail,
		AuthorURL:   paste.AuthorURL,

		Translate: hand.l10n.findLocale(req).translate,
	}

	// Get body line end
	switch lineend.GetLineEnd(paste.Body) {
	case "\r\n":
		tmplData.LineEnd = "CRLF"
	case "\r":
		tmplData.LineEnd = "CR"
	default:
		tmplData.LineEnd = "LF"
	}

	// Show paste
	return data.pastePage.Execute(rw, tmplData)
}
