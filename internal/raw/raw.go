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

package raw

import (
	"net/http"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"
)

type Data struct {
	log *logger.Logger
	db  *storage.DB
	cfg *config.Config
}

func Load(log *logger.Logger, db *storage.DB, cfg *config.Config) *Data {
	return &Data{
		log: log,
		db:  db,
		cfg: cfg,
	}
}

func (data *Data) Hand(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Server", model.UserAgent)

	err := data.rawHand(rw, req)

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
