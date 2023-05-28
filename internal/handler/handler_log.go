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
	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

func (hand *handler) logRequest(c *gin.Context, code int) {
	c.Set("request_logged", true)
	hand.log.Info(c.ClientIP(), c.Request.Method, code, c.Request.URL.Path, `"`+c.Request.UserAgent()+`"`)
}

func (hand *handler) logError(c *gin.Context, e model.Error) {
	c.Set("request_logged", true)
	hand.log.Error(c.ClientIP(), c.Request.Method, e.Code, c.Request.URL.Path, `"`+c.Request.UserAgent()+`",`, "error:", `"`+e.Error()+`"`)
}

func (hand *handler) writeErrorJSON(c *gin.Context, e error) {
	resp := model.ParseError(e)

	for key, val := range resp.Header {
		c.Header(key, val)
	}

	// Log request
	if resp.Code < 500 {
		hand.logRequest(c, resp.Code)
	} else {
		hand.logError(c, resp)
	}

	c.JSON(resp.Code, resp)
}

func (hand *handler) writeErrorPlain(c *gin.Context, e error) {
	resp := model.ParseError(e)

	for key, val := range resp.Header {
		c.Header(key, val)
	}

	// Log request
	if resp.Code < 500 {
		hand.logRequest(c, resp.Code)
	} else {
		hand.logError(c, resp)
	}

	c.Data(resp.Code, gin.MIMEPlain, []byte(resp.Text))
}
