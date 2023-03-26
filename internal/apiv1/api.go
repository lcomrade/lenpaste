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
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
)

type Data struct {
	log *logger.Logger
	db  *storage.DB
	cfg *config.Config

	lexers []string
}

func Load(log *logger.Logger, db *storage.DB, cfg *config.Config) *Data {
	lexers := chromaLexers.Names(false)

	return &Data{
		log:    log,
		db:     db,
		cfg:    cfg,
		lexers: lexers,
	}
}

func (data *Data) Hand(rw http.ResponseWriter, req *http.Request) {
	// Process request
	var err error

	rw.Header().Set("Server", model.Software+"/"+model.Version)

	switch req.URL.Path {
	case "/api/v1/new":
		err = data.newHand(rw, req)
	case "/api/v1/get":
		err = data.getHand(rw, req)
	case "/api/v1/getServerInfo":
		err = data.getServerInfoHand(rw, req)
	default:
		err = model.ErrNotFound
	}

	// Log
	if err == nil {
		data.log.HttpRequest(req, 200)

	} else {
		code, err := data.writeError(rw, req, err)
		if err != nil {
			data.log.HttpError(req, err)
		} else {
			data.log.HttpRequest(req, code)
		}
	}
}
