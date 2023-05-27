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

	"git.lcomrade.su/root/lenpaste/internal/lenpasswd"
	"git.lcomrade.su/root/lenpaste/internal/model"

	"github.com/gin-gonic/gin"
)

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

	AuthOk bool

	Translate func(string, ...interface{}) template.HTML
}

func (hand *handler) newPasteHand(c *gin.Context) {
	var err error

	// Check auth
	authOk := true

	if hand.cfg.Auth.Method == "lenpasswd" {
		authOk = false

		user, pass, authExist := req.BasicAuth()
		if authExist {
			authOk, err = lenpasswd.LoadAndCheck(hand.cfg.Paths.LenPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}

		if !authOk {
			rw.Header().Add("WWW-Authenticate", "Basic")
			rw.WriteHeader(401)
		}
	}

	// Create paste if need
	if req.Method == "POST" {
		pasteID, _, _, err := netshare.PasteAddFromForm(req, hand.db, hand.cfg, hand.lexers)
		if err != nil {
			return err
		}

		// Redirect to paste
		writeRedirect(rw, req, "/"+pasteID, 302)
		return nil
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
		AuthorDefault:      c.Cookie("author"),
		AuthorEmailDefault: c.Cookie("authorEmail"),
		AuthorURLDefault:   c.Cookie("authorURL"),
		AuthOk:             authOk,
		Translate:          hand.l10n.findLocale(req).translate,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")

	return hand.main.Execute(rw, tmplData)
}
