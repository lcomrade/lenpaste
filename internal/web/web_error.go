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
	"html/template"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

type errorTmpl struct {
	Code      int
	AdminName string
	AdminMail string
	Translate func(string, ...interface{}) template.HTML
}

func (hand *handler) writeError(c *gin.Context, e error) (int, error) {
	errData := errorTmpl{
		Code:      0,
		AdminName: hand.cfg.Public.AdminName,
		AdminMail: hand.cfg.Public.AdminMail,
		Translate: hand.l10n.findLocale(req).translate,
	}

	// Parse error
	resp := model.ParseError(e)

	if resp.Code > 499 {
		data.log.HttpError(req, e)
	}

	errData.Code = resp.Code

	// Prepare header
	c.Header("Content-Type", "text/plain; charset=utf-8")
	for key, val := range resp.Header {
		c.Header(key, val)
	}

	// Set HTTP status code
	rw.WriteHeader(resp.Code)

	// Send response body
	err := data.errorPage.Execute(rw, errData)
	if err != nil {
		return 500, err
	}

	return errData.Code, nil
}
