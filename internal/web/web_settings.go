package web

import (
	"html/template"
	"net/http"
)

type settingsTmpl struct {
	Language       string
	AvailableLangs []string
	Translate      func(string, ...interface{}) template.HTML
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

		langCookie, err := req.Cookie("lang")
		if err == nil {
			dataTmpl.Language = langCookie.Value
		}

		// Show page
		rw.Header().Set("Content-Type", "text/html")

		err = data.Settings.Execute(rw, dataTmpl)
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

		writeRedirect(rw, req, "/settings", 302)
	}
}
