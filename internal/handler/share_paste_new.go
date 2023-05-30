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
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lineend"
	"github.com/gin-gonic/gin"
)

func (hand *handler) pasteNew(c *gin.Context) (model.NewPaste, error) {
	// Check rate limit
	err := hand.db.RateLimitCheck(model.RLPasteNew, net.ParseIP(c.ClientIP()))
	if err != nil {
		return model.NewPaste{}, err
	}

	// Read form
	paste := model.Paste{
		Title:       c.PostForm("title"),
		Body:        c.PostForm("body"),
		Syntax:      c.PostForm("syntax"),
		DeleteTime:  0,
		OneUse:      false,
		Author:      c.PostForm("author"),
		AuthorEmail: c.PostForm("authorEmail"),
		AuthorURL:   c.PostForm("authorURL"),
	}

	// Remove new line from title
	paste.Title = strings.Replace(paste.Title, "\n", "", -1)
	paste.Title = strings.Replace(paste.Title, "\r", "", -1)
	paste.Title = strings.Replace(paste.Title, "\t", " ", -1)

	// Check title
	if utf8.RuneCountInString(paste.Title) > hand.cfg.Paste.TitleMaxLen && hand.cfg.Paste.TitleMaxLen > 0 {
		return model.NewPaste{}, model.ErrPayloadTooLarge
	}

	// Check paste body
	if paste.Body == "" {
		return model.NewPaste{}, model.ErrBadRequest
	}

	if utf8.RuneCountInString(paste.Body) > hand.cfg.Paste.BodyMaxLen && hand.cfg.Paste.BodyMaxLen > 0 {
		return model.NewPaste{}, model.ErrPayloadTooLarge
	}

	// Change paste body lines end
	switch c.PostForm("lineEnd") {
	case "", "LF", "lf":
		paste.Body = lineend.UnknownToUnix(paste.Body)

	case "CRLF", "crlf":
		paste.Body = lineend.UnknownToDos(paste.Body)

	case "CR", "cr":
		paste.Body = lineend.UnknownToOldMac(paste.Body)

	default:
		return model.NewPaste{}, model.ErrBadRequest
	}

	// Check syntax
	if paste.Syntax == "" {
		paste.Syntax = "plaintext"
	}

	syntaxOk := false
	for _, name := range hand.lexers {
		if name == paste.Syntax {
			syntaxOk = true
			break
		}
	}

	if !syntaxOk {
		return model.NewPaste{}, model.ErrBadRequest
	}

	// Get delete time
	expirStr := c.PostForm("expiration")
	if expirStr != "" {
		// Convert string to int
		expir, err := strconv.ParseInt(expirStr, 10, 64)
		if err != nil {
			return model.NewPaste{}, model.ErrBadRequest
		}

		// Check limits
		if hand.cfg.Paste.MaxLifetime > 0 {
			if expir > hand.cfg.Paste.MaxLifetime || expir <= 0 {
				return model.NewPaste{}, model.ErrBadRequest
			}
		}

		// Save if ok
		if expir > 0 {
			paste.DeleteTime = time.Now().Unix() + expir
		}
	}

	// Get "one use" parameter
	if c.PostForm("oneUse") == "true" {
		paste.OneUse = true
	}

	// Check author name, email and URL length.
	if utf8.RuneCountInString(paste.Author) > model.MaxLengthAuthorAll {
		return model.NewPaste{}, model.ErrPayloadTooLarge
	}

	if utf8.RuneCountInString(paste.AuthorEmail) > model.MaxLengthAuthorAll {
		return model.NewPaste{}, model.ErrPayloadTooLarge
	}

	if utf8.RuneCountInString(paste.AuthorURL) > model.MaxLengthAuthorAll {
		return model.NewPaste{}, model.ErrPayloadTooLarge
	}

	// Create paste
	data, err := hand.db.PasteAdd(paste)
	if err != nil {
		return model.NewPaste{}, err
	}

	return data, nil
}
