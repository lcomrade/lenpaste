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
)

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
		// If continue button not pressed
		req.ParseForm()

		if req.PostForm.Get("oneUseContinue") != "true" {
			err = data.PasteContinue.Execute(rw, paste)
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

	// Show paste
	err = data.PastePage.Execute(rw, paste)
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
