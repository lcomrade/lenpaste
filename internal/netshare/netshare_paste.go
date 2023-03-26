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

package netshare

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"git.lcomrade.su/root/lineend"
)

func PasteAddFromForm(req *http.Request, db *storage.DB, cfg *config.Config, lexerNames []string) (string, int64, int64, error) {
	// Check HTTP method
	if req.Method != "POST" {
		return "", 0, 0, model.ErrMethodNotAllowed
	}

	// Check rate limit
	err := db.RateLimitCheck("paste_new", GetClientAddr(req))
	if err != nil {
		return "", 0, 0, err
	}

	// Read form
	req.ParseForm()

	paste := model.Paste{
		Title:       req.PostForm.Get("title"),
		Body:        req.PostForm.Get("body"),
		Syntax:      req.PostForm.Get("syntax"),
		DeleteTime:  0,
		OneUse:      false,
		Author:      req.PostForm.Get("author"),
		AuthorEmail: req.PostForm.Get("authorEmail"),
		AuthorURL:   req.PostForm.Get("authorURL"),
	}

	// Remove new line from title
	paste.Title = strings.Replace(paste.Title, "\n", "", -1)
	paste.Title = strings.Replace(paste.Title, "\r", "", -1)
	paste.Title = strings.Replace(paste.Title, "\t", " ", -1)

	// Check title
	if utf8.RuneCountInString(paste.Title) > cfg.Paste.TitleMaxLen && cfg.Paste.TitleMaxLen >= 0 {
		return "", 0, 0, model.ErrPayloadTooLarge
	}

	// Check paste body
	if paste.Body == "" {
		return "", 0, 0, model.ErrBadRequest
	}

	if utf8.RuneCountInString(paste.Body) > cfg.Paste.BodyMaxLen && cfg.Paste.BodyMaxLen > 0 {
		return "", 0, 0, model.ErrPayloadTooLarge
	}

	// Change paste body lines end
	switch req.PostForm.Get("lineEnd") {
	case "", "LF", "lf":
		paste.Body = lineend.UnknownToUnix(paste.Body)

	case "CRLF", "crlf":
		paste.Body = lineend.UnknownToDos(paste.Body)

	case "CR", "cr":
		paste.Body = lineend.UnknownToOldMac(paste.Body)

	default:
		return "", 0, 0, model.ErrBadRequest
	}

	// Check syntax
	if paste.Syntax == "" {
		paste.Syntax = "plaintext"
	}

	syntaxOk := false
	for _, name := range lexerNames {
		if name == paste.Syntax {
			syntaxOk = true
			break
		}
	}

	if !syntaxOk {
		return "", 0, 0, model.ErrBadRequest
	}

	// Get delete time
	expirStr := req.PostForm.Get("expiration")
	if expirStr != "" {
		// Convert string to int
		expir, err := strconv.ParseInt(expirStr, 10, 64)
		if err != nil {
			return "", 0, 0, model.ErrBadRequest
		}

		// Check limits
		if cfg.Paste.MaxLifetime > 0 {
			if expir > cfg.Paste.MaxLifetime || expir <= 0 {
				return "", 0, 0, model.ErrBadRequest
			}
		}

		// Save if ok
		if expir > 0 {
			paste.DeleteTime = time.Now().Unix() + expir
		}
	}

	// Get "one use" parameter
	if req.PostForm.Get("oneUse") == "true" {
		paste.OneUse = true
	}

	// Check author name, email and URL length.
	if utf8.RuneCountInString(paste.Author) > model.MaxLengthAuthorAll {
		return "", 0, 0, model.ErrPayloadTooLarge
	}

	if utf8.RuneCountInString(paste.AuthorEmail) > model.MaxLengthAuthorAll {
		return "", 0, 0, model.ErrPayloadTooLarge
	}

	if utf8.RuneCountInString(paste.AuthorURL) > model.MaxLengthAuthorAll {
		return "", 0, 0, model.ErrPayloadTooLarge
	}

	// Create paste
	pasteID, createTime, deleteTime, err := db.PasteAdd(paste)
	if err != nil {
		return pasteID, createTime, deleteTime, err
	}

	return pasteID, createTime, deleteTime, nil
}
