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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Locale struct {
	Translate map[string]string `json:"translate"`
}

type Locales map[string]*Locale

func loadLocales(localeDir string) (Locales, error) {
	locales := make(Locales)

	// Get locale files list
	files, err := ioutil.ReadDir(localeDir)
	if err != nil {
		return locales, err
	}

	for _, fileInfo := range files {
		// Check file
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()
		if strings.HasSuffix(fileName, ".json") == false {
			continue
		}
		localeCode := fileName[:len(fileName)-5]

		// Read file
		file, err := os.Open(filepath.Join(localeDir, fileName))
		if err != nil {
			return locales, err
		}
		defer file.Close()

		var locale Locale
		err = json.NewDecoder(file).Decode(&locale)
		if err != nil {
			return locales, err
		}
		locales[localeCode] = &locale
	}

	return locales, nil
}

func (locales Locales) findLocale(req *http.Request) Locale {
	// Get accept language by cookie
	langCookie, err := req.Cookie("lang")
	if err == nil {
		if langCookie.Value != "" {
			locale, ok := locales[langCookie.Value]
			if ok != true {
				return *locale
			}
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
	for localeCode, locale := range locales {
		for _, lang := range langs {
			if lang == "en" {
				return Locale{}
			}

			if localeCode == lang {
				return *locale
			}
		}
	}

	return Locale{}
}

func (locale Locale) translate(s string) string {
	for val, key := range locale.Translate {
		if val == s {
			return key
		}
	}

	return s
}
