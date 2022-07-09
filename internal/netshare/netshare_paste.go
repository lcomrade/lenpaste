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
	"git.lcomrade.su/root/lineend"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func PasteAddFromForm(form url.Values, db storage.DB, titleMaxLen int, bodyMaxLen int, maxLifeTime int64, lexerNames []string) (string, error) {
	// Read form
	paste := storage.Paste{
		Title:      form.Get("title"),
		Body:       form.Get("body"),
		Syntax:     form.Get("syntax"),
		DeleteTime: 0,
		OneUse:     false,
	}

	// Remove new line from title
	paste.Title = strings.Replace(paste.Title, "\n", "", -1)
	paste.Title = strings.Replace(paste.Title, "\r", "", -1)

	// Check title
	if len(paste.Title) > titleMaxLen && titleMaxLen >= 0 {
		return "", ErrBadRequest
	}

	// Check paste body
	if paste.Body == "" {
		return "", ErrBadRequest
	}

	if len(paste.Body) > bodyMaxLen && bodyMaxLen > 0 {
		return "", ErrBadRequest
	}

	// Change paste body lines end
	switch form.Get("lineEnd") {
	case "", "LF", "lf":
		paste.Body = lineend.UnknownToUnix(paste.Body)

	case "CRLF", "crlf":
		paste.Body = lineend.UnknownToDos(paste.Body)

	case "CR", "cr":
		paste.Body = lineend.UnknownToOldMac(paste.Body)

	default:
		return "", ErrBadRequest
	}

	// Check syntax
	syntaxOk := false
	for _, name := range lexerNames {
		if name == paste.Syntax {
			syntaxOk = true
			break
		}
	}

	if syntaxOk == false {
		return "", ErrBadRequest
	}

	// Get delete time
	expirStr := form.Get("expiration")
	if expirStr != "" {
		// Convert string to int
		expir, err := strconv.ParseInt(expirStr, 10, 64)
		if err != nil {
			return "", ErrBadRequest
		}

		// Check limits
		if expir > maxLifeTime && maxLifeTime > 0 {
			return "", ErrBadRequest
		}

		// Save if ok
		if expir > 0 {
			paste.DeleteTime = time.Now().Unix() + expir
		}
	}

	// Get "one use" parameter
	if form.Get("oneUse") == "true" {
		paste.OneUse = true
	}

	// Create paste
	pasteID, err := db.PasteAdd(paste)
	if err != nil {
		return pasteID, err
	}

	return pasteID, nil
}
