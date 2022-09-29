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
	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/lexers"
	"html/template"
	"path/filepath"
	textTemplate "text/template"
)

type Data struct {
	DB  storage.DB
	Log logger.Config

	Lexers         *[]string
	Locales        *Locales
	LocaleSelector *map[string]string

	StyleCSS       *[]byte
	ErrorPage      *template.Template
	Main           *template.Template
	MainJS         *[]byte
	HistoryJS      *textTemplate.Template
	CodeJS         *textTemplate.Template
	PastePage      *template.Template
	PasteJS        *textTemplate.Template
	PasteContinue  *template.Template
	Settings       *template.Template
	About          *template.Template
	TermsOfUse     *template.Template
	Authors        *template.Template
	License        *template.Template
	SourceCodePage *template.Template

	Docs        *template.Template
	DocsApiV1   *template.Template
	DocsApiLibs *template.Template

	EmbeddedPage     *template.Template
	EmbeddedHelpPage *template.Template

	Version *string

	TitleMaxLen *int
	BodyMaxLen  *int
	MaxLifeTime *int64

	ServerAbout      *string
	ServerRules      *string
	ServerTermsExist *bool
	ServerTermsOfUse *string

	AdminName *string
	AdminMail *string

	RobotsDisallow *bool

	LenPasswdFile *string

	UiDefaultLifeTime *string
}

func Load(cfg config.Config, webDir string) (Data, error) {
	var data Data
	var err error

	// Setup base info
	data.DB = cfg.DB
	data.Log = cfg.Log

	data.Version = &cfg.Version

	data.TitleMaxLen = &cfg.TitleMaxLen
	data.BodyMaxLen = &cfg.BodyMaxLen
	data.MaxLifeTime = &cfg.MaxLifeTime
	data.UiDefaultLifeTime = &cfg.UiDefaultLifetime
	data.LenPasswdFile = &cfg.LenPasswdFile

	data.ServerAbout = &cfg.ServerAbout
	data.ServerRules = &cfg.ServerRules
	data.ServerTermsOfUse = &cfg.ServerTermsOfUse

	serverTermsExist := false
	if cfg.ServerTermsOfUse != "" {
		serverTermsExist = true
	}
	data.ServerTermsExist = &serverTermsExist

	data.AdminName = &cfg.AdminName
	data.AdminMail = &cfg.AdminMail

	data.RobotsDisallow = &cfg.RobotsDisallow

	// Get Chroma lexers
	lexers := chromaLexers.Names(false)
	data.Lexers = &lexers

	// Load locales
	locales, err := loadLocales(filepath.Join(webDir, "locale"))
	if err != nil {
		return data, err
	}
	data.Locales = &locales

	// Get locale selector
	localeSelector := make(map[string]string, len(locales))
	for key, val := range locales {
		valLoad := *val
		localeSelector[key] = valLoad["locale.Name"]
	}

	data.LocaleSelector = &localeSelector

	// style.css file
	styleCSS, err := readFile(filepath.Join(webDir, "style.css"))
	if err != nil {
		return data, err
	}
	data.StyleCSS = &styleCSS

	// main.tmpl
	data.Main, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "main.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// main.js
	mainJS, err := readFile(filepath.Join(webDir, "main.js"))
	if err != nil {
		return data, err
	}
	data.MainJS = &mainJS

	// history.js
	data.HistoryJS, err = textTemplate.ParseFiles(filepath.Join(webDir, "history.js"))
	if err != nil {
		return data, err
	}

	// code.js
	data.CodeJS, err = textTemplate.ParseFiles(filepath.Join(webDir, "code.js"))
	if err != nil {
		return data, err
	}

	// paste.tmpl
	data.PastePage, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "paste.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// paste.js
	data.PasteJS, err = textTemplate.ParseFiles(filepath.Join(webDir, "paste.js"))
	if err != nil {
		return data, err
	}

	// paste_continue.tmpl
	data.PasteContinue, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "paste_continue.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// settings.tmpl
	data.Settings, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "settings.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// about.tmpl
	data.About, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "about.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// terms.tmpl
	data.TermsOfUse, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "terms.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// authors.tmpl
	data.Authors, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "authors.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// license.tmpl
	data.License, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "license.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// source_code.tmpl
	data.SourceCodePage, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "source_code.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// docs.tmpl
	data.Docs, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "docs.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// docs_apiv1.tmpl
	data.DocsApiV1, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "docs_apiv1.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// docs_api_libs.tmpl
	data.DocsApiLibs, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "docs_api_libs.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// error.tmpl
	data.ErrorPage, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "error.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// emb.tmpl
	data.EmbeddedPage, err = template.ParseFiles(
		filepath.Join(webDir, "emb.tmpl"),
	)
	if err != nil {
		return data, err
	}

	// emb_help.tmpl
	data.EmbeddedHelpPage, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "emb_help.tmpl"),
	)
	if err != nil {
		return data, err
	}

	return data, nil
}
