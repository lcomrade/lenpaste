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

package apiv1

import (
	"encoding/json"
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/model"
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
func (data *Data) getServerInfoHand(rw http.ResponseWriter, req *http.Request) error {
	// Check method
	if req.Method != "GET" {
		return model.ErrMethodNotAllowed
	}

	// Get request parameters
	req.ParseForm()
	lang := req.Form.Get("lang")

	// Prepare data
	serverInfo := serverInfoType{
		Software:          model.Software,
		Version:           model.Version,
		TitleMaxLen:       data.cfg.Paste.TitleMaxLen,
		BodyMaxLen:        data.cfg.Paste.BodyMaxLen,
		MaxLifeTime:       data.cfg.Paste.MaxLifetime,
		ServerAbout:       data.cfg.GetAbout(lang),
		ServerRules:       data.cfg.GetRules(lang),
		ServerTermsOfUse:  data.cfg.GetTermsOfUse(lang),
		AdminName:         data.cfg.Public.AdminName,
		AdminMail:         data.cfg.Public.AdminMail,
		Syntaxes:          data.lexers,
		UiDefaultLifeTime: data.cfg.Paste.UiDefaultLifetimeStr,
		AuthRequired:      data.cfg.Paths.LenPasswdFile != "",
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(serverInfo)
}
