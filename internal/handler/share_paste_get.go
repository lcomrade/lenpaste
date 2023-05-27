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

package handler

import (
	"net"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

func (hand *handler) pasteGet(id string, openOneUse bool, clientIP net.IP) (model.Paste, error) {
	// Check rate limit
	err := hand.db.RateLimitCheck("paste_get", clientIP)
	if err != nil {
		return model.Paste{}, err
	}

	// Check paste id
	if id == "" {
		return model.Paste{}, model.ErrBadRequest
	}

	// Get paste
	paste, err := hand.db.PasteGet(id)
	if err != nil {
		return model.Paste{}, err
	}

	// If "one use" paste
	if paste.OneUse {
		if openOneUse {
			// Delete paste
			err = hand.db.PasteDelete(id)
			if err != nil {
				return model.Paste{}, err
			}

		} else {
			// Remove secret data
			paste = model.Paste{
				ID:     paste.ID,
				OneUse: true,
			}
		}
	}

	return paste, nil
}
