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
)

type Data struct {
	DB  storage.DB
	Log logger.Config

	Lexers *[]string

	StyleCSS       *[]byte
	RobotsTxt      *[]byte
	ErrorPage      *template.Template
	Main           *template.Template
	PastePage      *template.Template
	PasteContinue  *template.Template
	About          *template.Template
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

	ServerAbout *string
	ServerRules *string

	AdminName *string
	AdminMail *string
}

func Load(cfg config.Config, webDir string, robotsTxt []byte) (Data, error) {
	var data Data
	var err error

	// Setup base info
	data.DB = cfg.DB
	data.Log = cfg.Log

	data.Version = &cfg.Version

	data.TitleMaxLen = &cfg.TitleMaxLen
	data.BodyMaxLen = &cfg.BodyMaxLen
	data.MaxLifeTime = &cfg.MaxLifeTime

	data.ServerAbout = &cfg.ServerAbout
	data.ServerRules = &cfg.ServerRules

	data.AdminName = &cfg.AdminName
	data.AdminMail = &cfg.AdminMail

	data.RobotsTxt = &robotsTxt

	// Get Chroma lexers
	lexers := chromaLexers.Names(false)
	data.Lexers = &lexers

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

	// paste.tmpl
	data.PastePage, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "paste.tmpl"),
	)
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

	// about.tmpl
	data.About, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "about.tmpl"),
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
