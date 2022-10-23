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
	"html/template"
	"net/http"
	"time"
)

type embTmpl struct {
	ID            string
	CreateTimeStr string
	DeleteTime    int64
	OneUse        bool
	Title         string
	Body          template.HTML

	ErrorNotFound bool
	Translate     func(string, ...interface{}) template.HTML
}

// Pattern: /emb/
func (data Data) EmbeddedHand(rw http.ResponseWriter, req *http.Request) {
	errorNotFound := false

	// Log request
	data.Log.HttpRequest(req)

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[5:])

	// Read DB
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		if err == storage.ErrNotFoundID {
			errorNotFound = true

		} else {
			data.writeError(rw, req, err)
			return
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()

	tmplData := embTmpl{
		ID:            paste.ID,
		CreateTimeStr: createTime.Format("1 Jan, 2006"),
		DeleteTime:    paste.DeleteTime,
		OneUse:        paste.OneUse,
		Title:         paste.Title,
		Body:          tryHighlight(paste.Body, paste.Syntax),

		ErrorNotFound: errorNotFound,
		Translate:     data.Locales.findLocale(req).translate,
	}

	// Show paste
	err = data.EmbeddedPage.Execute(rw, tmplData)
	if err != nil {
		data.writeError(rw, req, err)
		return
	}
}
