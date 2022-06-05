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

package web

import (
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"net/http"
)

type embHelpTmpl struct {
	ID       string
	Protocol string
	Host     string
}

// Pattern: /emb_help/
func (data Data) EmbeddedHelpHand(rw http.ResponseWriter, req *http.Request) {
	// Log request
	data.Log.HttpRequest(req)

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[10:])

	// Read DB
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		if err == storage.ErrNotFoundID {
			data.errorNotFound(rw, req)
			return

		} else {
			data.errorInternal(rw, req, err)
			return
		}
	}

	// Show paste
	tmplData := embHelpTmpl{
		ID:       paste.ID,
		Protocol: netshare.GetProtocol(req.Header),
		Host:     netshare.GetHost(req),
	}

	err = data.EmbeddedHelpPage.Execute(rw, tmplData)
	if err != nil {
		data.errorInternal(rw, req, err)
		return
	}
}
