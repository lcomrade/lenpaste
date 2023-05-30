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
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/gin-gonic/gin"
)

type serverInfoType struct {
	Software          string   `json:"software"`
	Version           string   `json:"version"`
	TitleMaxLen       int      `json:"titleMaxlength"`
	BodyMaxLen        int      `json:"bodyMaxlength"`
	MaxLifeTime       int64    `json:"maxLifeTime"`
	ServerAbout       string   `json:"serverAbout"`
	ServerRules       string   `json:"serverRules"`
	ServerTermsOfUse  string   `json:"serverTermsOfUse"`
	AdminName         string   `json:"adminName"`
	AdminMail         string   `json:"adminMail"`
	Syntaxes          []string `json:"syntaxes"`
	UiDefaultLifeTime string   `json:"uiDefaultLifeTime"`
	AuthRequired      bool     `json:"authRequired"`
}

// GET /api/v1/getServerInfo
func (hand *handler) getServerInfoHand(c *gin.Context) {
	// Get request parameters
	lang := c.Query("lang")

	// Prepare data
	serverInfo := serverInfoType{
		Software:          model.Software,
		Version:           model.Version,
		TitleMaxLen:       hand.cfg.Paste.TitleMaxLen,
		BodyMaxLen:        hand.cfg.Paste.BodyMaxLen,
		MaxLifeTime:       hand.cfg.Paste.MaxLifetime,
		ServerAbout:       hand.cfg.GetAbout(lang),
		ServerRules:       hand.cfg.GetRules(lang),
		ServerTermsOfUse:  hand.cfg.GetTermsOfUse(lang),
		AdminName:         hand.cfg.Public.AdminName,
		AdminMail:         hand.cfg.Public.AdminMail,
		Syntaxes:          hand.lexers,
		UiDefaultLifeTime: hand.cfg.Paste.UiDefaultLifetime,
		AuthRequired:      false,
	}

	// Return response
	c.JSON(http.StatusOK, serverInfo)
}
