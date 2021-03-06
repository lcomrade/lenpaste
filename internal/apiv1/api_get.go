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
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"net/http"
)

// GET /api/v1/get
func (data Data) GetHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Check method
	if req.Method != "GET" {
		data.writeError(rw, req, netshare.ErrMethodNotAllowed)
		return
	}

	// Get paste ID
	req.ParseForm()

	pasteID := req.Form.Get("id")

	// Check paste id
	if pasteID == "" {
		data.writeError(rw, req, netshare.ErrBadRequest)
		return
	}

	// Get paste
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		data.writeError(rw, req, err)
		return
	}

	// If "one use" paste
	if paste.OneUse == true {
		if req.Form.Get("openOneUse") == "true" {
			// Delete paste
			err = data.DB.PasteDelete(pasteID)
			if err != nil {
				data.writeError(rw, req, err)
				return
			}

		} else {
			// Remove secret data
			paste = storage.Paste{
				ID:     paste.ID,
				OneUse: true,
			}
		}
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(rw).Encode(paste)
	if err != nil {
		data.Log.HttpError(req, err)
		return
	}
}
