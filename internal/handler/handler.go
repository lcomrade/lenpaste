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
	"embed"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"

	"html/template"
	textTemplate "text/template"
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

	spaceJS  *textTemplate.Template
	spaceCSS *textTemplate.Template
}

func Run(log *logger.Logger, db *storage.DB, cfg *config.Config) error {
	var err error
	hand := &handler{
		log:    log,
		db:     db,
		cfg:    cfg,
		lexers: chromaLexers.Names(false),
	}

	// Setup Gin logging
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Debug("[GIN-debug]", httpMethod, absolutePath, "-->", handlerName, "("+strconv.Itoa(nuHandlers)+" handlers)")
	}

	if !model.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure HTTP router
	r := gin.New()

	if len(cfg.HTTP.TrustedProxies) != 0 {
		r.ForwardedByClientIP = true
		err := r.SetTrustedProxies(cfg.HTTP.TrustedProxies)
		if err != nil {
			return errors.New("handler: run: " + err.Error())
		}

	} else {
		log.Warning("You trusted all proxies, this is NOT safe. Change \"http.trusted_proxies\" value in config.")
	}

	r.Use(func(c *gin.Context) {
		c.Next()

		// Check request logging
		if c.GetBool("request_logged") {
			return
		}

		c.Set("request_logged", true)

		// Log request if need
		hand.logRequest(c, http.StatusOK)
	})

	r.Use(gin.Recovery())

	// Setup HTTP paths
	r.GET("/api/v1/get", hand.apiPasteGet)
	r.POST("/api/v1/new", hand.apiPasteNew)

	r.GET("/api/v1/getServerInfo", hand.getServerInfoHand)

	r.GET("/raw/:id", hand.rawHand)

	// Load locales
	hand.l10n, err = loadLocales(embFS, "data/locale")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// Load themes
	hand.themes, err = loadThemes(cfg.Paths.ThemesDir, hand.l10n, cfg.UI.DefaultTheme)
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// style.css file
	hand.styleCSS, err = textTemplate.ParseFS(embFS, "data/style.css")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// main.js
	mainJS, err := embFS.ReadFile("data/main.js")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}
	hand.mainJS = &mainJS

	// history.js
	hand.historyJS, err = textTemplate.ParseFS(embFS, "data/history.js")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// code.js
	hand.codeJS, err = textTemplate.ParseFS(embFS, "data/code.js")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// paste.js
	hand.pasteJS, err = textTemplate.ParseFS(embFS, "data/paste.js")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// space.css
	hand.spaceCSS, err = textTemplate.ParseFS(embFS, "data/space.css")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// space.js
	hand.spaceJS, err = textTemplate.ParseFS(embFS, "data/space.js")
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	// Load templates
	{
		d, err := embFS.ReadDir("data")
		if err != nil {
			return errors.New("handler: " + err.Error())
		}

		// Add some html/template files manual
		tmplFilesManual := [][]string{
			{"emb.tmpl"},
			{"space.tmpl"},
		}

		// Add other html/template files
		var tmplFiles []string
		for _, file := range d {
			if file.IsDir() {
				continue
			}

			// If file already manual added
			fileName := file.Name()

			ok := true
			for _, part := range tmplFilesManual {
				if fileName == part[len(part)-1] {
					ok = false
					break
				}
			}

			if !ok {
				continue
			}

			// Else save file name
			if strings.HasSuffix(fileName, ".tmpl") {
				tmplFiles = append(tmplFiles, fileName)
			}
		}

		// Load templates to Gin
		render := multitemplate.NewRenderer()

		for _, fileName := range tmplFiles {
			render.Add(fileName, template.Must(
				template.ParseFS(embFS,
					"data/base.tmpl",
					"data/"+fileName,
				),
			))
		}

		for _, part := range tmplFilesManual {
			var tmpls []string
			for _, fileName := range part {
				tmpls = append(tmpls, "data/"+fileName)
			}

			render.Add(part[len(part)-1], template.Must(
				template.ParseFS(embFS, tmpls...),
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
	r.POST("/", hand.newPasteHand)
	r.GET("/:id", hand.getPasteHand)

	r.GET("/space/", hand.spaceHand)
	r.GET("/space/:id", hand.spaceHand)
	r.GET("/space.css", hand.spaceCSSHand)
	r.GET("/space.js", hand.spaceJSHand)

	r.GET("/settings", hand.settingsHand)
	r.POST("/settings", hand.settingsHand)
	r.GET("/terms", hand.termsOfUseHand)

	r.GET("/dl/:id", hand.dlHand)
	r.GET("/emb/:id", hand.embeddedHand)
	r.GET("/emb_help/:id", hand.embeddedHelpHand)

	// Run
	err = r.Run(cfg.HTTP.Address)
	if err != nil {
		return errors.New("handler: " + err.Error())
	}

	return nil
}
