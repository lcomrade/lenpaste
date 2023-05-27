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

	"github.com/gin-gonic/gin"
)

func (hand *handler) apiPasteGet(c *gin.Context) {
	openOneUse := false
	if c.Query("openOneUse") == "true" {
		openOneUse = true
	}

	paste, err := hand.pasteGet(c.Query("id"), openOneUse, net.ParseIP(c.ClientIP()))
	if err != nil {
		hand.writeErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, paste)
}

func (hand *handler) apiPasteNew(c *gin.Context) {
	newPaste, err := hand.pasteNew(c)
	if err != nil {
		hand.writeErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, newPaste)
}
