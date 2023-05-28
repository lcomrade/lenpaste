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

func (hand *handler) newPasteHand(c *gin.Context) {
	type createTmpl struct {
		TitleMaxLen       int
		BodyMaxLen        int
		AuthorAllMaxLen   int
		MaxLifeTime       int64
		UiDefaultLifeTime string
		Lexers            []string
		ServerTermsExist  bool

		AuthorDefault      string
		AuthorEmailDefault string
		AuthorURLDefault   string

		Translate func(string, ...interface{}) template.HTML
	}

	// Create paste if need
	if c.Request.Method == "POST" {
		newPaste, err := hand.pasteNew(c)
		if err != nil {
			hand.writeErrorWeb(c, err)
			return
		}

		// Redirect to paste
		c.Redirect(http.StatusFound, "/"+newPaste.ID)
		return
	}

	// Else show create page
	tmplData := createTmpl{
		TitleMaxLen:        hand.cfg.Paste.TitleMaxLen,
		BodyMaxLen:         hand.cfg.Paste.BodyMaxLen,
		AuthorAllMaxLen:    model.MaxLengthAuthorAll,
		MaxLifeTime:        hand.cfg.Paste.MaxLifetime,
		UiDefaultLifeTime:  hand.cfg.Paste.UiDefaultLifetime,
		Lexers:             hand.lexers,
		ServerTermsExist:   hand.cfg.TermsOfUse != nil,
		AuthorDefault:      getCookie(c, "author"),
		AuthorEmailDefault: getCookie(c, "authorEmail"),
		AuthorURLDefault:   getCookie(c, "authorURL"),
		Translate:          hand.l10n.findLocale(c).translate,
	}

	c.HTML(http.StatusOK, "main.tmpl", tmplData)
}
