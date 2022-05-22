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

import(
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"net/http"
	"time"
)

// Pattern: /new
func (data Data) CreatePasteHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Read request
	req.ParseForm()

	if req.PostForm.Get("body") != "" {
		paste := storage.Paste{
			Title: req.PostForm.Get("title"),
			Body: req.PostForm.Get("body"),
			OneUse: false,
		}
	
		// Get expiration time
		expirTime, err := storage.ExpirationToTime(req.PostForm.Get("expiration"))
		if err != nil {
			data.errorBadRequest(rw, req)
			return
		}
		paste.DeleteTime = time.Now().Unix() + expirTime
	
		// Get "one use" parameter
		if req.PostForm.Get("oneUse") == "true" {
			paste.OneUse = true
		}

		// Create paste
		_, err = data.DB.PasteAdd(paste)
		if err != nil {
			data.errorInternal(rw, req, err)
			return
		}

		return
	}

	// Else show create page
	rw.Header().Set("Content-Type", "text/html")

	err := data.CreatePaste.Execute(rw, "")
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
