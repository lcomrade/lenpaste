package web

import (
	"html/template"
	"net/http"
)

type termsOfUseTmpl struct {
	TermsOfUse string

	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /terms
func (data Data) TermsOfUseHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Show page
	rw.Header().Set("Content-Type", "text/html")

	err := data.TermsOfUse.Execute(rw, termsOfUseTmpl{
		TermsOfUse: *data.ServerTermsOfUse,
		Highlight:  tryHighlight,
		Translate:  data.Locales.findLocale(req).translate},
	)
	if err != nil {
		data.writeError(rw, req, err)
	}
}
