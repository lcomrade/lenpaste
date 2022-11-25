package web

import (
	"git.lcomrade.su/root/lenpaste/internal/lenpasswd"
	"html/template"
	"net/http"
)

const cookieMaxAge = 60 * 60 * 24 * 360 * 50 // 50 year

type settingsTmpl struct {
	Language         string
	LanguageSelector map[string]string

	Author      string
	AuthorEmail string
	AuthorURL   string

	AuthOk bool

	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /settings
func (data *Data) SettingsHand(rw http.ResponseWriter, req *http.Request) {
	var err error

	// Log request
	data.Log.HttpRequest(req)

	// Check auth
	authOk := true

	if *data.LenPasswdFile != "" {
		authOk = false

		user, pass, authExist := req.BasicAuth()
		if authExist == true {
			authOk, err = lenpasswd.LoadAndCheck(*data.LenPasswdFile, user, pass)
			if err != nil {
				data.writeError(rw, req, err)
				return
			}
		}
	}

	// Show settings page
	if req.Method != "POST" {
		// Prepare data
		dataTmpl := settingsTmpl{
			Language:         getCookie(req, "lang"),
			LanguageSelector: data.LocalesList,
			Author:           getCookie(req, "author"),
			AuthorEmail:      getCookie(req, "authorEmail"),
			AuthorURL:        getCookie(req, "authorURL"),
			AuthOk:           authOk,
			Translate:        data.Locales.findLocale(req).translate,
		}

		// Show page
		rw.Header().Set("Content-Type", "text/html")

		err := data.Settings.Execute(rw, dataTmpl)
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
}
