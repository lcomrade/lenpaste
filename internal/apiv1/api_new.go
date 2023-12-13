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
	"github.com/lcomrade/lenpaste/internal/lenpasswd"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"net/http"
)

type newPasteAnswer struct {
	ID         string `json:"id"`
	CreateTime int64  `json:"createTime"`
	DeleteTime int64  `json:"deleteTime"`
}

// POST /api/v1/new
func (data *Data) newHand(rw http.ResponseWriter, req *http.Request) error {
	var err error

	// Check auth
	if data.LenPasswdFile != "" {
		authOk := false

		user, pass, authExist := req.BasicAuth()
		if authExist == true {
			authOk, err = lenpasswd.LoadAndCheck(data.LenPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}

		if authOk == false {
			return netshare.ErrUnauthorized
		}
	}

	// Check method
	if req.Method != "POST" {
		return netshare.ErrMethodNotAllowed
	}

	// Get form data and create paste
	pasteID, createTime, deleteTime, err := netshare.PasteAddFromForm(req, data.DB, data.RateLimitNew, data.TitleMaxLen, data.BodyMaxLen, data.MaxLifeTime, data.Lexers)
	if err != nil {
		return err
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(newPasteAnswer{ID: pasteID, CreateTime: createTime, DeleteTime: deleteTime})
}
