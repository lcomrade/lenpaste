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
	"net/http"
	"strings"
	textTemplate "text/template"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
)

//go:embed data/*
var embFS embed.FS

type Data struct {
	log *logger.Logger
	db  *storage.DB
	cfg *config.Config

	lexers []string
	l10n   *l10n
	themes *themesData

	styleCSS       *textTemplate.Template
	errorPage      *template.Template
	main           *template.Template
	mainJS         *[]byte
	historyJS      *textTemplate.Template
	codeJS         *textTemplate.Template
	pastePage      *template.Template
	pasteJS        *textTemplate.Template
	pasteContinue  *template.Template
	settings       *template.Template
	about          *template.Template
	termsOfUse     *template.Template
	authors        *template.Template
	license        *template.Template
	sourceCodePage *template.Template

	docs        *template.Template
	docsApiV1   *template.Template
	docsApiV2   *template.Template
	docsApiLibs *template.Template

	embeddedPage     *template.Template
	embeddedHelpPage *template.Template
}

func Load(log *logger.Logger, db *storage.DB, cfg *config.Config) (*Data, error) {
	// Setup base info
	var err error
	data := Data{
		log: log,
		db:  db,
		cfg: cfg,
	}

	// Get Chroma lexers
	data.lexers = chromaLexers.Names(false)

	// Load locales
	data.l10n, err = loadLocales(embFS, "data/locale")
	if err != nil {
		return nil, err
	}

	// Load themes
	data.themes, err = loadThemes(cfg.ThemesDir, data.l10n, cfg.UI.DefaultTheme)
	if err != nil {
		return nil, err
	}

	// style.css file
	data.styleCSS, err = textTemplate.ParseFS(embFS, "data/style.css")
	if err != nil {
		return nil, err
	}

	// main.tmpl
	data.main, err = template.ParseFS(embFS, "data/base.tmpl", "data/main.tmpl")
	if err != nil {
		return nil, err
	}

	// main.js
	mainJS, err := embFS.ReadFile("data/main.js")
	if err != nil {
		return nil, err
	}
	data.mainJS = &mainJS

	// history.js
	data.historyJS, err = textTemplate.ParseFS(embFS, "data/history.js")
	if err != nil {
		return nil, err
	}

	// code.js
	data.codeJS, err = textTemplate.ParseFS(embFS, "data/code.js")
	if err != nil {
		return nil, err
	}

	// paste.tmpl
	data.pastePage, err = template.ParseFS(embFS, "data/base.tmpl", "data/paste.tmpl")
	if err != nil {
		return nil, err
	}

	// paste.js
	data.pasteJS, err = textTemplate.ParseFS(embFS, "data/paste.js")
	if err != nil {
		return nil, err
	}

	// paste_continue.tmpl
	data.pasteContinue, err = template.ParseFS(embFS, "data/base.tmpl", "data/paste_continue.tmpl")
	if err != nil {
		return nil, err
	}

	// settings.tmpl
	data.settings, err = template.ParseFS(embFS, "data/base.tmpl", "data/settings.tmpl")
	if err != nil {
		return nil, err
	}

	// about.tmpl
	data.about, err = template.ParseFS(embFS, "data/base.tmpl", "data/about.tmpl")
	if err != nil {
		return nil, err
	}

	// terms.tmpl
	data.termsOfUse, err = template.ParseFS(embFS, "data/base.tmpl", "data/terms.tmpl")
	if err != nil {
		return nil, err
	}

	// authors.tmpl
	data.authors, err = template.ParseFS(embFS, "data/base.tmpl", "data/authors.tmpl")
	if err != nil {
		return nil, err
	}

	// license.tmpl
	data.license, err = template.ParseFS(embFS, "data/base.tmpl", "data/license.tmpl")
	if err != nil {
		return nil, err
	}

	// source_code.tmpl
	data.sourceCodePage, err = template.ParseFS(embFS, "data/base.tmpl", "data/source_code.tmpl")
	if err != nil {
		return nil, err
	}

	// docs.tmpl
	data.docs, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_apiv1.tmpl
	data.docsApiV1, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_apiv1.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_apiv2.tmpl
	data.docsApiV2, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_apiv2.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_api_libs.tmpl
	data.docsApiLibs, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_api_libs.tmpl")
	if err != nil {
		return nil, err
	}

	// error.tmpl
	data.errorPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/error.tmpl")
	if err != nil {
		return nil, err
	}

	// emb.tmpl
	data.embeddedPage, err = template.ParseFS(embFS, "data/emb.tmpl")
	if err != nil {
		return nil, err
	}

	// emb_help.tmpl
	data.embeddedHelpPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/emb_help.tmpl")
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (data *Data) Handler(rw http.ResponseWriter, req *http.Request) {
	// Process request
	var err error

	rw.Header().Set("Server", model.UserAgent)

	switch req.URL.Path {
	// Search engines
	case "/robots.txt":
		err = data.robotsTxtHand(rw, req)
	case "/sitemap.xml":
		err = data.sitemapHand(rw, req)
	// Resources
	case "/style.css":
		err = data.styleCSSHand(rw, req)
	case "/main.js":
		err = data.mainJSHand(rw, req)
	case "/history.js":
		err = data.historyJSHand(rw, req)
	case "/code.js":
		err = data.codeJSHand(rw, req)
	case "/paste.js":
		err = data.pasteJSHand(rw, req)
	case "/about":
		err = data.aboutHand(rw, req)
	case "/about/authors":
		err = data.authorsHand(rw, req)
	case "/about/license":
		err = data.licenseHand(rw, req)
	case "/about/source_code":
		err = data.sourceCodePageHand(rw, req)
	case "/docs":
		err = data.docsHand(rw, req)
	case "/docs/apiv1":
		err = data.docsApiV1Hand(rw, req)
	case "/docs/apiv2":
		err = data.docsApiV2Hand(rw, req)
	case "/docs/api_libs":
		err = data.docsApiLibsHand(rw, req)
	// Pages
	case "/":
		err = data.newPasteHand(rw, req)
	case "/settings":
		err = data.settingsHand(rw, req)
	case "/terms":
		err = data.termsOfUseHand(rw, req)
	// Else
	default:
		if strings.HasPrefix(req.URL.Path, "/dl/") {
			err = data.dlHand(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/emb/") {
			err = data.embeddedHand(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/emb_help/") {
			err = data.embeddedHelpHand(rw, req)

		} else {
			err = data.getPasteHand(rw, req)
		}
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
