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
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/lenpasswd"
	"git.lcomrade.su/root/lenpaste/internal/model"
)

const cookieMaxAge = 60 * 60 * 24 * 360 * 50 // 50 year

type settingsTmpl struct {
	Language         string
	LanguageSelector map[string]string

	Theme         string
	ThemeSelector map[string]string

	AuthorAllMaxLen int
	Author          string
	AuthorEmail     string
	AuthorURL       string

	AuthOk bool

	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /settings
func (data *Data) settingsHand(rw http.ResponseWriter, req *http.Request) error {
	var err error

	// Check auth
	authOk := true

	if data.cfg.LenPasswdFile != "" {
		authOk = false

		user, pass, authExist := req.BasicAuth()
		if authExist {
			authOk, err = lenpasswd.LoadAndCheck(data.cfg.LenPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}
	}

	// Show settings page
	if req.Method != "POST" {
		// Prepare data
		dataTmpl := settingsTmpl{
			Language:         getCookie(req, "lang"),
			LanguageSelector: data.l10n.names,
			Theme:            getCookie(req, "theme"),
			ThemeSelector:    data.themes.getForLocale(data.l10n, req),
			AuthorAllMaxLen:  model.MaxLengthAuthorAll,
			Author:           getCookie(req, "author"),
			AuthorEmail:      getCookie(req, "authorEmail"),
			AuthorURL:        getCookie(req, "authorURL"),
			AuthOk:           authOk,
			Translate:        data.l10n.findLocale(req).translate,
		}

		if dataTmpl.Theme == "" {
			dataTmpl.Theme = data.cfg.UI.DefaultTheme
		}

		// Show page
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := data.settings.Execute(rw, dataTmpl)
		if err != nil {
			data.writeError(rw, req, err)
		}

		// Else update settings
	} else {
		req.ParseForm()

		lang := req.PostForm.Get("lang")
		if lang == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   "lang",
				Value:  "",
				MaxAge: -1,
			})

		} else {
			http.SetCookie(rw, &http.Cookie{
				Name:   "lang",
				Value:  lang,
				MaxAge: cookieMaxAge,
			})
		}

		theme := req.PostForm.Get("theme")
		if theme == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   "theme",
				Value:  "",
				MaxAge: -1,
			})

		} else {
			http.SetCookie(rw, &http.Cookie{
				Name:   "theme",
				Value:  theme,
				MaxAge: cookieMaxAge,
			})
		}

		author := req.PostForm.Get("author")
		if author == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   "author",
				Value:  "",
				MaxAge: -1,
			})

		} else {
			http.SetCookie(rw, &http.Cookie{
				Name:   "author",
				Value:  author,
				MaxAge: cookieMaxAge,
			})
		}

		authorEmail := req.PostForm.Get("authorEmail")
		if authorEmail == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   "authorEmail",
				Value:  "",
				MaxAge: -1,
			})

		} else {
			http.SetCookie(rw, &http.Cookie{
				Name:   "authorEmail",
				Value:  authorEmail,
				MaxAge: cookieMaxAge,
			})
		}

		authorURL := req.PostForm.Get("authorURL")
		if authorURL == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   "authorURL",
				Value:  "",
				MaxAge: -1,
			})

		} else {
			http.SetCookie(rw, &http.Cookie{
				Name:   "authorURL",
				Value:  authorURL,
				MaxAge: cookieMaxAge,
			})
		}

		writeRedirect(rw, req, "/settings", 302)
	}

	return nil
}
