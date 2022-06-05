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
	Title         string
	Body          template.HTML
}

// Pattern: /emb/
func (data Data) EmbeddedHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[5:])

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
		// Return 404 error
		data.errorNotFound(rw, req)
		return
	}

	// Highlight body
	bodyHighlight := tryHighlight(paste.Body, paste.Syntax)

	// Prepare template dat
	createTime := time.Unix(paste.CreateTime, 0).UTC()

	tmplData := embTmpl{
		ID:            paste.ID,
		CreateTimeStr: createTime.Format("1 Jan, 2006"),
		Title:         paste.Title,
		Body:          template.HTML(bodyHighlight),
	}

	// Show paste
	err = data.EmbeddedPage.Execute(rw, tmplData)
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
