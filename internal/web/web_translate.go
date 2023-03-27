// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHdata ANY WARRANTY; withdata even the implied warranty of MERCHANTABILITY
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

	"git.lcomrade.su/root/lenpaste/internal/model"
)

type locale map[string]string

type l10n struct {
	// VALUE = locales[LANG_CODE][KEY]
	//
	// For example:
	//   locales["en"]["base.About"] = "About"
	//   locales["ru"]["base.About"] = "О сайте"
	locales map[string]locale

	// LANG_CODE - PRETTY_NAME
	//
	// For example:
	//   en - English
	//   ru - Русский
	names map[string]string
}

func loadLocales(f embed.FS, localeDir string) (*l10n, error) {
	data := l10n{
		locales: make(map[string]locale),
		names:   make(map[string]string),
	}

	// Get locale files list
	files, err := f.ReadDir(localeDir)
	if err != nil {
		return nil, errors.New("web: failed read dir \"" + localeDir + "\": " + err.Error())
	}

	// Load locales
	for _, fileInfo := range files {
		// Check file
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()
		if !strings.HasSuffix(fileName, ".json") {
			continue
		}
		localeCode := fileName[:len(fileName)-5]

		// Open and read file
		filePath := filepath.Join(localeDir, fileName)
		file, err := f.Open(filePath)
		if err != nil {
			return nil, errors.New("web: failed open file \"" + filePath + "\": " + err.Error())
		}
		defer file.Close()

		var locale map[string]string
		err = json.NewDecoder(file).Decode(&locale)
		if err != nil {
			return nil, errors.New("web: failed read file \"" + filePath + "\": " + err.Error())
		}

		data.locales[localeCode] = locale
	}

	// Prepare locales names
	for key, val := range data.locales {
		// Get locale name
		localeName := val["locale.Name"]
		if localeName == "" {
			return nil, errors.New("web: empty locale.Name parameter in \"" + key + "\" locale")
		}

		// Append to the translation, if it is not complete
		defLocale := data.locales[model.BaseLocale]
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
			return nil, errors.New("web: locale \"" + key + "\" is empty")
		}

		if curTotal == defTotal {
			data.names[key] = localeName
		} else {
			data.names[key] = localeName + fmt.Sprintf(" (%.2f%%)", (float32(curTotal)/float32(defTotal))*100)
		}
	}

	return &data, nil
}

func (data *l10n) detectLanguage(req *http.Request) string {
	// Get accept language from cookie
	{
		lang := getCookie(req, "lang")
		if lang != "" {
			_, ok := data.locales[lang]
			if ok {
				return lang
			}
		}
	}

	// Get user Accept-Languages list
	acceptLanguage := req.Header.Get("Accept-Language")
	acceptLanguage = strings.Replace(acceptLanguage, " ", "", -1)

	var langs []string
	for _, part := range strings.Split(acceptLanguage, ";") {
		for _, lang := range strings.Split(part, ",") {
			if !strings.HasPrefix(lang, "q=") {
				langs = append(langs, lang)
			}
		}
	}

	// Search locale
	for _, lang := range langs {
		for localeCode := range data.locales {
			if localeCode == lang {
				return lang
			}
		}
	}

	// Else return default language
	return model.BaseLocale
}

func (data *l10n) findLocale(req *http.Request) locale {
	return data.locales[data.detectLanguage(req)]
}

func (locale locale) translate(s string, a ...interface{}) template.HTML {
	for key, val := range locale {
		if key == s {
			return template.HTML(fmt.Sprintf(val, a...))
		}
	}

	panic(errors.New("web: translate: unknown locale key \"" + s + "\""))
}
