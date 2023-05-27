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

package web

import (
	"net/http"
	"strings"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/netshare"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/gin-gonic/gin"
)

// Pattern: /dl/
func (hand *handler) dlHand(c *gin.Context) {
	// Check rate limit
	err := hand.db.RateLimitCheck("paste_get", netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Read DB
	pasteID := string([]rune(req.URL.Path)[4:])

	paste, err := hand.db.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse {
		// Delete paste
		err = hand.db.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Get create time
	createTime := time.Unix(paste.CreateTime, 0).UTC()

	// Get file name
	fileName := paste.ID
	if paste.Title != "" {
		fileName = paste.Title
	}

	// Get file extension
	fileExt := chromaLexers.Get(paste.Syntax).Config().Filenames[0][1:]
	if !strings.HasSuffix(fileName, fileExt) {
		fileName = fileName + fileExt
	}

	// Write result
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	http.ServeContent(rw, req, fileName, createTime, strings.NewReader(paste.Body))

	return nil
}
