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

package raw

import (
	"github.com/lcomrade/lenpaste/internal/netshare"
	"io"
	"net/http"
)

// Pattern: /raw/
func (data *Data) rawHand(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Read DB
	pasteID := string([]rune(req.URL.Path)[5:])

	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse == true {
		// Delete paste
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Write result
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err = io.WriteString(rw, paste.Body)
	if err != nil {
		return err
	}

	return nil
}
