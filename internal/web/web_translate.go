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
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type Locale map[string]string
type Locales map[string]*Locale

func loadLocales(f embed.FS, localeDir string) (Locales, error) {
	locales := make(Locales)

	// Get locale files list
	files, err := f.ReadDir(localeDir)
	if err != nil {
		return locales, err
	}

	for _, fileInfo := range files {
		// Check file
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()
		if strings.HasSuffix(fileName, ".locale") == false {
			continue
		}
		localeCode := fileName[:len(fileName)-7]

		// Read file
		filePath := filepath.Join(localeDir, fileName)
		fileByte, err := f.ReadFile(filePath)
		if err != nil {
			return locales, err
		}

		fileStr := bytes.NewBuffer(fileByte).String()

		// Load locale
		locale, err := readKVCfg(fileStr)
		if err != nil {
			return locales, errors.New("web: failed read file '" + filePath + "': " + err.Error())
		}

		localeLocale := Locale(locale)
		locales[localeCode] = &localeLocale
	}

	return locales, nil
}

func (locales Locales) findLocale(req *http.Request) *Locale {
	// Get accept language by cookie
	langCookie := getCookie(req, "lang")
	if langCookie != "" {
		locale, ok := locales[langCookie]
		if ok == true {
			return locale
		}
	}

	// Get user Accepr-Languages list
	acceptLanguage := req.Header.Get("Accept-Language")
	acceptLanguage = strings.Replace(acceptLanguage, " ", "", -1)

	var langs []string
	for _, part := range strings.Split(acceptLanguage, ";") {
		for _, lang := range strings.Split(part, ",") {
			if strings.HasPrefix(lang, "q=") == false {
				langs = append(langs, lang)
			}
		}
	}

	// Search locale
	for _, lang := range langs {
		for localeCode, locale := range locales {
			if localeCode == lang {
				return locale
			}
		}
	}

	// Load default locale
	locale, ok := locales["en"]
	if ok != true {
		// If en locale not found load first locale
		for _, l := range locales {
			return l
		}
	}

	return locale
}

func (locale *Locale) translate(s string, a ...interface{}) template.HTML {
	for key, val := range *locale {
		if key == s {
			return template.HTML(fmt.Sprintf(val, a...))
		}
	}

	return template.HTML(s)
}
