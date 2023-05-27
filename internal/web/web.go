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

package web

import (
	"embed"
	"html/template"
	"strings"
	textTemplate "text/template"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

//go:embed data/*
var embFS embed.FS

type handler struct {
	log *logger.Logger
	db  *storage.DB
	cfg *config.Config

	lexers []string
	l10n   *l10n
	themes *themesData

	styleCSS  *textTemplate.Template
	mainJS    *[]byte
	historyJS *textTemplate.Template
	codeJS    *textTemplate.Template
	pasteJS   *textTemplate.Template
}

func Install(log *logger.Logger, db *storage.DB, cfg *config.Config, r *gin.Engine) error {
	// Setup base info
	var err error
	hand := &handler{
		log: log,
		db:  db,
		cfg: cfg,
	}

	// Get Chroma lexers
	hand.lexers = chromaLexers.Names(false)

	// Load locales
	hand.l10n, err = loadLocales(embFS, "data/locale")
	if err != nil {
		return err
	}

	// Load themes
	hand.themes, err = loadThemes(cfg.Paths.ThemesDir, hand.l10n, cfg.UI.DefaultTheme)
	if err != nil {
		return err
	}

	// style.css file
	hand.styleCSS, err = textTemplate.ParseFS(embFS, "data/style.css")
	if err != nil {
		return err
	}

	// main.js
	mainJS, err := embFS.ReadFile("data/main.js")
	if err != nil {
		return err
	}
	hand.mainJS = &mainJS

	// history.js
	hand.historyJS, err = textTemplate.ParseFS(embFS, "data/history.js")
	if err != nil {
		return err
	}

	// code.js
	hand.codeJS, err = textTemplate.ParseFS(embFS, "data/code.js")
	if err != nil {
		return err
	}

	// paste.js
	hand.pasteJS, err = textTemplate.ParseFS(embFS, "data/paste.js")
	if err != nil {
		return err
	}

	// Load templates
	{
		d, err := embFS.ReadDir("data/")
		if err != nil {
			return err
		}

		var tmplFiles []string
		for _, file := range d {
			fileName := file.Name()

			if strings.HasSuffix(fileName, ".tmpl") {
				tmplFiles = append(tmplFiles, fileName[:len(fileName)-5])
			}
		}

		render := multitemplate.NewRenderer()

		for _, fileName := range tmplFiles {
			render.Add(fileName, template.Must(
				template.ParseFS(embFS,
					"data/base.tmpl",
					"data/"+fileName+".tmpl",
				),
			))
		}

		r.HTMLRender = render
	}

	// Setup paths
	r.GET("/robots.txt", hand.robotsTxtHand)
	r.GET("/sitemap.xml", hand.sitemapHand)

	r.GET("/style.css", hand.styleCSSHand)
	r.GET("/main.js", hand.mainJSHand)
	r.GET("/history.js", hand.historyJSHand)
	r.GET("/code.js", hand.codeJSHand)
	r.GET("/paste.js", hand.pasteJSHand)

	r.GET("/about", hand.aboutHand)
	r.GET("/about/authors", hand.authorsHand)
	r.GET("/about/license", hand.licenseHand)
	r.GET("/about/source_code", hand.sourceCodePageHand)

	r.GET("/docs", hand.docsHand)
	r.GET("/docs/apiv1", hand.docsApiV1Hand)
	r.GET("/docs/api_libs", hand.docsApiLibsHand)

	r.GET("/", hand.newPasteHand)
	r.GET("/:id", hand.getPasteHand)
	r.GET("/settings", hand.settingsHand)
	r.GET("/terms", hand.termsOfUseHand)

	r.GET("/dl/", hand.dlHand)
	r.GET("/emb/", hand.embeddedHand)
	r.GET("/emb_help/", hand.embeddedHelpHand)

	return nil
}
