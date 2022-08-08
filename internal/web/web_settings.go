package web

import (
	"html/template"
	"net/http"
)

type settingsTmpl struct {
	Language       string
	AvailableLangs []string

	Author      string
	AuthorEmail string
	AuthorURL   string

	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /settings
func (data Data) SettingsHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Show settings page
	if req.Method != "POST" {
		// Prepare data
		var dataTmpl settingsTmpl
		dataTmpl.Translate = data.Locales.findLocale(req).translate

		for code, _ := range *data.Locales {
			dataTmpl.AvailableLangs = append(dataTmpl.AvailableLangs, code)
		}

		dataTmpl.Language = getCookie(req, "lang")
		dataTmpl.Author = getCookie(req, "author")
		dataTmpl.AuthorEmail = getCookie(req, "authorEmail")
		dataTmpl.AuthorURL = getCookie(req, "authorURL")

		// Show page
		rw.Header().Set("Content-Type", "text/html")

		err := data.Settings.Execute(rw, dataTmpl)
		if err != nil {
			data.errorInternal(rw, req, err)
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
				MaxAge: 0,
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
				MaxAge: 0,
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
				MaxAge: 0,
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
				MaxAge: 0,
			})
		}

		writeRedirect(rw, req, "/settings", 302)
	}
}
