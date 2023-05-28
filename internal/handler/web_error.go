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
	"html/template"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

func (hand *handler) writeErrorWeb(c *gin.Context, e error) {
	type errorTmpl struct {
		Code      int
		AdminName string
		AdminMail string
		Translate func(string, ...interface{}) template.HTML
	}

	errData := errorTmpl{
		Code:      0,
		AdminName: hand.cfg.Public.AdminName,
		AdminMail: hand.cfg.Public.AdminMail,
		Translate: hand.l10n.findLocale(c).translate,
	}

	// Parse error
	resp := model.ParseError(e)

	if resp.Code > 499 {
		hand.logError(c, resp)
	}

	errData.Code = resp.Code

	// Prepare header
	for key, val := range resp.Header {
		c.Header(key, val)
	}

	// Send response body
	c.HTML(resp.Code, "error.tmpl", errData)
}
