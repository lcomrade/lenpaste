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

package apiv1

import (
	"encoding/json"
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"net/http"
)

type newPasteAnswer struct {
	ID string `json:"id"`
}

// POST /api/v1/new
func (data Data) NewHand(rw http.ResponseWriter, req *http.Request) {
	// Check method
	if req.Method != "POST" {
		data.writeError(rw, req, netshare.ErrMethodNotAllowed)
		return
	}

	// Get form data and create paste
	req.ParseForm()

	pasteID, err := netshare.PasteAddFromForm(req.PostForm, data.DB, *data.TitleMaxLen, *data.BodyMaxLen, *data.MaxLifeTime, *data.Lexers)
	if err != nil {
		if err == netshare.ErrBadRequest {
			data.writeError(rw, req, netshare.ErrBadRequest)
			return
		}

		data.writeError(rw, req, netshare.ErrInternal)
		return
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(rw).Encode(newPasteAnswer{ID: pasteID})
	if err != nil {
		data.Log.HttpError(req, err)
		return
	}
}
