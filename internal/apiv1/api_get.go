// Copyright (C) 2021-2023 Leonid Maslakov.

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
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/netshare"
)

// GET /api/v1/get
func (data *Data) getHand(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.db.RateLimitCheck("paste_get", netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Check method
	if req.Method != "GET" {
		return model.ErrMethodNotAllowed
	}

	// Get paste ID
	req.ParseForm()

	pasteID := req.Form.Get("id")

	// Check paste id
	if pasteID == "" {
		return model.ErrBadRequest
	}

	// Get paste
	paste, err := data.db.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse {
		if req.Form.Get("openOneUse") == "true" {
			// Delete paste
			err = data.db.PasteDelete(pasteID)
			if err != nil {
				return err
			}

		} else {
			// Remove secret data
			paste = model.Paste{
				ID:     paste.ID,
				OneUse: true,
			}
		}
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(paste)
}
