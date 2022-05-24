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

package storage

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrFormBadRequest = errors.New("Bad Request")
	ErrInternal       = errors.New("Internal Server Error")
)

func (dbInfo DB) PasteAddFromForm(form url.Values) error {
	// Read form
	paste := Paste{
		Title:      form.Get("title"),
		Body:       form.Get("body"),
		DeleteTime: 0,
		OneUse:     false,
	}

	// Get delete time
	expirStr := form.Get("expiration")
	if expirStr != "" {
		expir, err := strconv.ParseInt(expirStr, 16, 64)
		if err != nil {
			return ErrFormBadRequest
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
		return ErrInternal
	}

	return nil
}
