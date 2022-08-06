// Copyright (C) 2021-2022 Leonid Maslakov.

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
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"git.lcomrade.su/root/lineend"
	"html/template"
	"net/http"
	"time"
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
	Translate     func(string, ...interface{}) template.HTML
}

type pasteContinueTmpl struct {
	ID        string
	Translate func(string, ...interface{}) template.HTML
}

func (data Data) getPaste(rw http.ResponseWriter, req *http.Request) {
	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[1:])

	// Read DB
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		if err == storage.ErrNotFoundID {
			data.errorNotFound(rw, req)
			return

		} else {
			data.errorInternal(rw, req, err)
			return
		}
	}

	// If "one use" paste
	if paste.OneUse == true {
		// If continue button not pressed
		req.ParseForm()

		if req.PostForm.Get("oneUseContinue") != "true" {
			tmplData := pasteContinueTmpl{
				ID:        paste.ID,
				Translate: data.Locales.findLocale(req).translate,
			}

			err = data.PasteContinue.Execute(rw, tmplData)
			if err != nil {
				data.errorInternal(rw, req, err)
				return
			}

			return
		}

		// If continue button pressed delete paste
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			data.errorInternal(rw, req, err)
			return
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()
	deleteTime := time.Unix(paste.DeleteTime, 0).UTC()

	tmplData := pasteTmpl{
		ID:         paste.ID,
		Title:      paste.Title,
		Body:       tryHighlight(paste.Body, paste.Syntax),
		Syntax:     paste.Syntax,
		CreateTime: paste.CreateTime,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,

		CreateTimeStr: createTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		DeleteTimeStr: deleteTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),

		Translate: data.Locales.findLocale(req).translate,
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
	err = data.PastePage.Execute(rw, tmplData)
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
