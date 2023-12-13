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
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"net/http"
	"strings"
	"time"
)

// Pattern: /dl/
func (data *Data) dlHand(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Read DB
	pasteID := string([]rune(req.URL.Path)[4:])

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

	// Get create time
	createTime := time.Unix(paste.CreateTime, 0).UTC()

	// Get file name
	fileName := paste.ID
	if paste.Title != "" {
		fileName = paste.Title
	}

	// Get file extension
	fileExt := chromaLexers.Get(paste.Syntax).Config().Filenames[0][1:]
	if strings.HasSuffix(fileName, fileExt) == false {
		fileName = fileName + fileExt
	}

	// Write result
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	rw.Header().Set("Content-Transfer-Encoding", "binary")
	rw.Header().Set("Expires", "0")

	http.ServeContent(rw, req, fileName, createTime, strings.NewReader(paste.Body))

	return nil
}
