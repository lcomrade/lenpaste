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
	"github.com/lcomrade/lenpaste/internal/lenpasswd"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"html/template"
	"net/http"
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

func (data *Data) newPasteHand(rw http.ResponseWriter, req *http.Request) error {
	var err error

	// Check auth
	authOk := true

	if data.LenPasswdFile != "" {
		authOk = false

		user, pass, authExist := req.BasicAuth()
		if authExist == true {
			authOk, err = lenpasswd.LoadAndCheck(data.LenPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}

		if authOk == false {
			rw.Header().Add("WWW-Authenticate", "Basic")
			rw.WriteHeader(401)
		}
	}

	// Create paste if need
	if req.Method == "POST" {
		pasteID, _, _, err := netshare.PasteAddFromForm(req, data.DB, data.RateLimitNew, data.TitleMaxLen, data.BodyMaxLen, data.MaxLifeTime, data.Lexers)
		if err != nil {
			return err
		}

		// Redirect to paste
		writeRedirect(rw, req, "/"+pasteID, 302)
		return nil
	}

	// Else show create page
	tmplData := createTmpl{
		TitleMaxLen:        data.TitleMaxLen,
		BodyMaxLen:         data.BodyMaxLen,
		AuthorAllMaxLen:    netshare.MaxLengthAuthorAll,
		MaxLifeTime:        data.MaxLifeTime,
		UiDefaultLifeTime:  data.UiDefaultLifeTime,
		Lexers:             data.Lexers,
		ServerTermsExist:   data.ServerTermsExist,
		AuthorDefault:      getCookie(req, "author"),
		AuthorEmailDefault: getCookie(req, "authorEmail"),
		AuthorURLDefault:   getCookie(req, "authorURL"),
		AuthOk:             authOk,
		Translate:          data.Locales.findLocale(req).translate,
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	return data.Main.Execute(rw, tmplData)
}
