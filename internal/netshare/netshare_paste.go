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

package netshare

import (
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"net/url"
	"strconv"
	"time"
)

func PasteAddFromForm(dbInfo storage.DB, form url.Values) (storage.Paste, error) {
	// Read form
	paste := storage.Paste{
		Title:      form.Get("title"),
		Body:       form.Get("body"),
		DeleteTime: 0,
		OneUse:     false,
	}

	// Check paste body
	if paste.Body == "" {
		return paste, ErrBadRequest
	}

	// Get delete time
	expirStr := form.Get("expiration")
	if expirStr != "" {
		expir, err := strconv.ParseInt(expirStr, 10, 64)
		if err != nil {
			return paste, ErrBadRequest
		}

		if expir > 0 {
			paste.DeleteTime = time.Now().Unix() + expir
		}
	}

	// Get "one use" parameter
	if form.Get("oneUse") == "true" {
		paste.OneUse = true
	}

	// Create paste
	paste, err := dbInfo.PasteAdd(paste)
	if err != nil {
		return paste, err
	}

	return paste, nil
}
