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
	"errors"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/gin-gonic/gin"
)

type handler struct {
	log *logger.Logger
	db  *storage.DB
	cfg *config.Config

	lexers []string
}

func Run(log *logger.Logger, db *storage.DB, cfg *config.Config) error {
	hand := &handler{
		log:    log,
		db:     db,
		cfg:    cfg,
		lexers: chromaLexers.Names(false),
	}

	// Setup Gin logging
	if model.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup Gin routers
	r := gin.New()

	r.GET("/api/v1/get", hand.apiPasteGet)
	r.POST("/api/v1/new", hand.apiPasteNew)

	r.GET("/api/v1/pasteGet", hand.apiPasteGet)
	r.PUT("/api/v1/pasteNew", hand.apiPasteNew)

	r.GET("/api/v1/getServerInfo", hand.getServerInfoHand)

	r.GET("/raw/:id", hand.rawHand)

	// Run
	err := r.Run(cfg.HTTP.Address)
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	return nil
}
