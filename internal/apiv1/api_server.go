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
	"github.com/lcomrade/lenpaste/internal/netshare"
	"net/http"
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
		return netshare.ErrMethodNotAllowed
	}

	// Prepare data
	serverInfo := serverInfoType{
		Software:          "Lenpaste",
		Version:           data.Version,
		TitleMaxLen:       data.TitleMaxLen,
		BodyMaxLen:        data.BodyMaxLen,
		MaxLifeTime:       data.MaxLifeTime,
		ServerAbout:       data.ServerAbout,
		ServerRules:       data.ServerRules,
		ServerTermsOfUse:  data.ServerTermsOfUse,
		AdminName:         data.AdminName,
		AdminMail:         data.AdminMail,
		Syntaxes:          data.Lexers,
		UiDefaultLifeTime: data.UiDefaultLifeTime,
		AuthRequired:      data.LenPasswdFile != "",
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(serverInfo)
}
