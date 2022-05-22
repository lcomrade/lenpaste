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
	"net/http"
	"strconv"
	"time"
)

// Pattern: /
func (data Data) MainHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// If main page
	if req.URL.Path == "/" {
		data.newPaste(rw, req)

		// Else show paste
	} else {
		data.getPaste(rw, req)
	}
}

func (data Data) newPaste(rw http.ResponseWriter, req *http.Request) {
	// Read request
	req.ParseForm()

	if req.PostForm.Get("body") != "" {
		paste := storage.Paste{
			Title:  req.PostForm.Get("title"),
			Body:   req.PostForm.Get("body"),
			OneUse: false,
		}

		// Get delete time
		expirStr := req.Form.Get("expiration")
		if expirStr != "" {
			expir, err := strconv.ParseInt(expirStr, 16, 64)
			if err != nil {
				data.errorBadRequest(rw, req)
				return
			}

			if expir > 0 {
				paste.DeleteTime = time.Now().Unix() + expir
			}
		}

		// Get "one use" parameter
		if req.PostForm.Get("oneUse") == "true" {
			paste.OneUse = true
		}

		// Create paste
		paste, err := data.DB.PasteAdd(paste)
		if err != nil {
			data.errorInternal(rw, req, err)
			return
		}

		// Redirect to paste
		writeRedirect(rw, req, "/"+paste.ID, 302)
		return
	}

	// Else show create page
	rw.Header().Set("Content-Type", "text/html")

	err := data.Main.Execute(rw, "")
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}

func (data Data) getPaste(rw http.ResponseWriter, req *http.Request) {
	// Read DB
	pasteID := string([]rune(req.URL.Path)[1:])

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
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			data.errorInternal(rw, req, err)
			return
		}
	}

	// Show paste
	err = data.PastePage.Execute(rw, paste)
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
