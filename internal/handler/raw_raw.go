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
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

// Pattern: /raw/
func (hand *handler) rawHand(c *gin.Context) {
	// Check rate limit
	err := hand.db.RateLimitCheck(model.RLPasteGet, net.IP(c.ClientIP()))
	if err != nil {
		hand.writeErrorPlain(c, err)
		return
	}

	// Read DB
	pasteID := c.Param("id")

	paste, err := hand.db.PasteGet(pasteID)
	if err != nil {
		hand.writeErrorPlain(c, err)
		return
	}

	// If "one use" paste
	if paste.OneUse {
		// Delete paste
		err = hand.db.PasteDelete(pasteID)
		if err != nil {
			hand.writeErrorPlain(c, err)
			return
		}
	}

	// Write result
	c.Data(http.StatusOK, gin.MIMEPlain, []byte(paste.Body))
}
