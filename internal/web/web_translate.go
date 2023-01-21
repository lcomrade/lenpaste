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
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

const baseLocale = "en"

type Locale map[string]string
type Locales map[string]Locale
type LocalesList map[string]string

func loadLocales(f embed.FS, localeDir string) (Locales, LocalesList, error) {
	locales := make(Locales)
	localesList := make(LocalesList)

	// Get locale files list
	files, err := f.ReadDir(localeDir)
	if err != nil {
		return nil, nil, errors.New("web: failed read dir '" + localeDir + "': " + err.Error())
	}

	// Load locales
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

		// Open and read file
		filePath := filepath.Join(localeDir, fileName)
		file, err := f.Open(filePath)
		if err != nil {
			return nil, nil, errors.New("web: failed open file '" + filePath + "': " + err.Error())
		}
		defer file.Close()

		var locale Locale
		err = json.NewDecoder(file).Decode(&locale)
		if err != nil {
			return nil, nil, errors.New("web: failed read file '" + filePath + "': " + err.Error())
		}

		locales[localeCode] = Locale(locale)
	}

	// Prepare locales list
	for key, val := range locales {
		// Get locale name
		localeName := val["locale.Name"]
		if localeName == "" {
			return nil, nil, errors.New("web: empty locale.Name parameter in '" + key + "' locale")
		}

		// Append to the translation, if it is not complete
		defLocale := locales[baseLocale]
		defTotal := len(defLocale)
		curTotal := 0
		for defKey, defVal := range defLocale {
			_, isExist := val[defKey]
			if isExist {
				curTotal = curTotal + 1
			} else {
				val[defKey] = defVal
			}
		}

		if curTotal == 0 {
			return nil, nil, errors.New("web: locale '" + key + "' is empty")
		}

		if curTotal == defTotal {
			localesList[key] = localeName
		} else {
			localesList[key] = localeName + fmt.Sprintf(" (%.2f%%)", (float32(curTotal)/float32(defTotal))*100)
		}
	}

	return locales, localesList, nil
}

func (locales Locales) findLocale(req *http.Request) Locale {
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
	locale, _ := locales[baseLocale]
	return locale
}

func (locale Locale) translate(s string, a ...interface{}) template.HTML {
	for key, val := range locale {
		if key == s {
			return template.HTML(fmt.Sprintf(val, a...))
		}
	}

	panic(errors.New("web: translate: unknown locale key: " + s))
}
