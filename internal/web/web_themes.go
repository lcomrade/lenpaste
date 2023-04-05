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
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

const embThemesDir = "data/theme"

type theme map[string]string

type themesData struct {
	// VALUE = themes[THEME_NAME][KEY]
	//
	// For example:
	//   locales["dark"]["color.Font"]  = "#FFFFFF"
	//   locales["light"]["color.Font"] = "#000000"
	themes map[string]theme

	// LOCALIZED_NAME = names[LANG_CODE][THEME_NAME]
	//
	// For example:
	//   names["en"]["dark"] = "Dark"
	//   names["ru"]["dark"] = "Тёмная"
	names map[string]map[string]string
}

func loadThemes(hostThemeDir string, locale *l10n, defaultTheme string) (*themesData, error) {
	themes := themesData{
		themes: make(map[string]theme),
		names:  make(map[string]map[string]string),
	}

	for localeCode := range locale.locales {
		themes.names[localeCode] = make(map[string]string)
	}

	// Prepare load FS function
	loadThemesFromFS := func(f fs.FS, themeDir string, ignoreIfNotExist bool) error {
		// Get theme files list
		files, err := fs.ReadDir(f, themeDir)
		if err != nil {
			if ignoreIfNotExist && os.IsNotExist(err) {
				return nil
			}
			return errors.New("web: failed read dir \"" + themeDir + "\": " + err.Error())
		}

		for _, fileInfo := range files {
			// Check file
			if fileInfo.IsDir() {
				continue
			}

			fileName := fileInfo.Name()
			if !strings.HasSuffix(fileName, ".theme") {
				continue
			}
			themeCode := fileName[:len(fileName)-6]

			// Read file
			filePath := filepath.Join(themeDir, fileName)
			fileByte, err := fs.ReadFile(f, filePath)
			if err != nil {
				return errors.New("web: failed open file \"" + filePath + "\": " + err.Error())
			}

			fileStr := bytes.NewBuffer(fileByte).String()

			// Load theme
			theme, err := readKVCfg(fileStr)
			if err != nil {
				return errors.New("web: failed read file \"" + filePath + "\": " + err.Error())
			}

			_, themeExist := themes.themes[themeCode]
			if themeExist {
				return errors.New("web: theme already loaded: " + filePath)
			}

			themes.themes[themeCode] = theme
		}

		return nil
	}

	// Load embed themes
	err := loadThemesFromFS(embFS, embThemesDir, false)
	if err != nil {
		return nil, err
	}

	// Load external themes
	err = loadThemesFromFS(os.DirFS(hostThemeDir), ".", true)
	if err != nil {
		return nil, err
	}

	// Prepare themes names list
	for key, val := range themes.themes {
		// Get theme name
		themeName := val["theme.Name."+model.BaseLocale]
		if themeName == "" {
			return nil, errors.New("web: empty \"theme.Name." + model.BaseLocale + "\" parameter in \"" + key + "\" theme")
		}

		// Append to the translation, if it is not complete
		defTheme := themes.themes[model.BaseTheme]
		defTotal := len(defTheme)
		curTotal := 0
		for defKey, defVal := range defTheme {
			_, isExist := val[defKey]
			if isExist {
				curTotal = curTotal + 1
			} else {
				if strings.HasPrefix(defKey, "theme.Name.") {
					val[defKey] = val["theme.Name."+model.BaseLocale]
				} else {
					val[defKey] = defVal
				}
			}
		}

		if curTotal == 0 {
			return nil, errors.New("web: theme \"" + key + "\" is empty")
		}

		// Add theme to themes names list
		themeNameSuffix := ""
		if curTotal != defTotal {
			themeNameSuffix = fmt.Sprintf(" (%.2f%%)", (float32(curTotal)/float32(defTotal))*100)
		}
		themes.names[model.BaseLocale][key] = themeName + themeNameSuffix

		for localeCode := range locale.locales {
			result, ok := val["theme.Name."+localeCode]
			if ok {
				themes.names[localeCode][key] = result + themeNameSuffix
			} else {
				themes.names[localeCode][key] = themeName + themeNameSuffix
			}
		}
	}

	// Check default theme exist
	_, ok := themes.themes[defaultTheme]
	if !ok {
		return nil, errors.New("web: default theme '" + defaultTheme + "' not found")
	}

	return &themes, nil
}

func (themes *themesData) getForLocale(locale *l10n, req *http.Request) map[string]string {
	return themes.names[locale.detectLanguage(req)]
}

func (data *themesData) findTheme(req *http.Request, defaultTheme string) theme {
	// Get theme by cookie
	themeCookie := getCookie(req, "theme")
	if themeCookie != "" {
		theme, ok := data.themes[themeCookie]
		if ok {
			return theme
		}
	}

	// Load default theme
	return data.themes[defaultTheme]
}

func (theme theme) theme(s string) string {
	for key, val := range theme {
		if key == s {
			return val
		}
	}

	panic(errors.New("web: theme: unknown theme key \"" + s + "\""))
}

func (theme theme) tryHighlight(source string, lexer string) template.HTML {
	return tryHighlight(source, lexer, theme.theme("highlight.Theme"))
}
