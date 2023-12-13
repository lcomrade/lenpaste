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
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/lcomrade/lenpaste/internal/config"
	"github.com/lcomrade/lenpaste/internal/logger"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"github.com/lcomrade/lenpaste/internal/storage"
	"net/http"
)

type Data struct {
	Log logger.Logger
	DB  storage.DB

	RateLimitNew *netshare.RateLimitSystem
	RateLimitGet *netshare.RateLimitSystem

	Lexers []string

	Version string

	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsOfUse string

	AdminName string
	AdminMail string

	LenPasswdFile string

	UiDefaultLifeTime string
}

func Load(db storage.DB, cfg config.Config) *Data {
	lexers := chromaLexers.Names(false)

	return &Data{
		DB:                db,
		Log:               cfg.Log,
		RateLimitNew:      cfg.RateLimitNew,
		RateLimitGet:      cfg.RateLimitGet,
		Lexers:            lexers,
		Version:           cfg.Version,
		TitleMaxLen:       cfg.TitleMaxLen,
		BodyMaxLen:        cfg.BodyMaxLen,
		MaxLifeTime:       cfg.MaxLifeTime,
		ServerAbout:       cfg.ServerAbout,
		ServerRules:       cfg.ServerRules,
		ServerTermsOfUse:  cfg.ServerTermsOfUse,
		AdminName:         cfg.AdminName,
		AdminMail:         cfg.AdminMail,
		LenPasswdFile:     cfg.LenPasswdFile,
		UiDefaultLifeTime: cfg.UiDefaultLifetime,
	}
}

func (data *Data) Hand(rw http.ResponseWriter, req *http.Request) {
	// Process request
	var err error

	rw.Header().Set("Server", config.Software+"/"+data.Version)

	switch req.URL.Path {
	// Search engines
	case "/api/v1/new":
		err = data.newHand(rw, req)
	case "/api/v1/get":
		err = data.getHand(rw, req)
	case "/api/v1/getServerInfo":
		err = data.getServerInfoHand(rw, req)
	default:
		err = netshare.ErrNotFound
	}

	// Log
	if err == nil {
		data.Log.HttpRequest(req, 200)

	} else {
		code, err := data.writeError(rw, req, err)
		if err != nil {
			data.Log.HttpError(req, err)
		} else {
			data.Log.HttpRequest(req, code)
		}
	}
}
