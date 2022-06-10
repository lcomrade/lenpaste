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

package apiv1

import (
	"encoding/json"
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"net/http"
)

type serverInfoType struct {
	Version     string   `json:"version"`
	TitleMaxLen int      `json:"titleMaxlength"`
	BodyMaxLen  int      `json:"bodyMaxlength"`
	Syntaxes    []string `json:"syntaxes"`
}

// GET /api/v1/getServerInfo
func (data Data) GetServerInfoHand(rw http.ResponseWriter, req *http.Request) {
	// Check method
	if req.Method != "GET" {
		data.writeError(rw, req, netshare.ErrBadRequest)
		return
	}

	// Prepare data
	serverInfo := serverInfoType{
		TitleMaxLen: data.TitleMaxLen,
		BodyMaxLen:  data.BodyMaxLen,
		Version:     data.Version,
		Syntaxes:    data.Lexers,
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(serverInfo)
	if err != nil {
		data.Log.HttpError(req, err)
		return
	}
}
