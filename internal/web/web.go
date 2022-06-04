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
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/lexers"
	"html/template"
	"path/filepath"
)

type Data struct {
	DB  storage.DB
	Log logger.Config

	Lexers []string

	StyleCSS       []byte
	Main           *template.Template
	PastePage      *template.Template
	PasteContinue  *template.Template
	About          *template.Template
	License        *template.Template
	SourceCodePage *template.Template
	Docs           *template.Template
	DocsApiV1      *template.Template
	ErrorPage      *template.Template

	Version string
}

func Load(webDir string, db storage.DB, log logger.Config, version string) (Data, error) {
	var data Data
	var err error

	// Setup base info
	data.DB = db
	data.Log = log
	data.Version = version

	// Get Chroma lexers
	data.Lexers = chromaLexers.Names(false)

	// style.css file
	data.StyleCSS, err = readFile(filepath.Join(webDir, "style.css"))
	if err != nil {
		return data, err
	}

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

	// error.tmpl
	data.ErrorPage, err = template.ParseFiles(
		filepath.Join(webDir, "base.tmpl"),
		filepath.Join(webDir, "error.tmpl"),
	)
	if err != nil {
		return data, err
	}

	return data, nil
}
