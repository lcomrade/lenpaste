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
)

const baseTheme = "dark"
const embThemesDir = "data/theme"

type Theme map[string]string
type Themes map[string]Theme

type ThemesListPart map[string]string
type ThemesList map[string]ThemesListPart

func loadThemes(hostThemeDir string, localesList LocalesList, defaultTheme string) (Themes, ThemesList, error) {
	themes := make(Themes)
	themesList := make(ThemesList)

	for localeCode, _ := range localesList {
		themesList[localeCode] = make(ThemesListPart)
	}

	// Prepare load FS function
	loadThemesFromFS := func(f fs.FS, themeDir string) error {
		// Get theme files list
		files, err := fs.ReadDir(f, themeDir)
		if err != nil {
			return errors.New("web: failed read dir '" + themeDir + "': " + err.Error())
		}

		for _, fileInfo := range files {
			// Check file
			if fileInfo.IsDir() {
				continue
			}

			fileName := fileInfo.Name()
			if strings.HasSuffix(fileName, ".theme") == false {
				continue
			}
			themeCode := fileName[:len(fileName)-6]

			// Read file
			filePath := filepath.Join(themeDir, fileName)
			fileByte, err := fs.ReadFile(f, filePath)
			if err != nil {
				return errors.New("web: failed open file '" + filePath + "': " + err.Error())
			}

			fileStr := bytes.NewBuffer(fileByte).String()

			// Load theme
			theme, err := readKVCfg(fileStr)
			if err != nil {
				return errors.New("web: failed read file '" + filePath + "': " + err.Error())
			}

			_, themeExist := themes[themeCode]
			if themeExist {
				return errors.New("web: theme alredy loaded: " + filePath)
			}

			themes[themeCode] = Theme(theme)
		}

		return nil
	}

	// Load embed themes
	err := loadThemesFromFS(embFS, embThemesDir)
	if err != nil {
		return nil, nil, err
	}

	// Load external themes
	if hostThemeDir != "" {
		err = loadThemesFromFS(os.DirFS(hostThemeDir), ".")
		if err != nil {
			return nil, nil, err
		}
	}

	// Prepare themes list
	for key, val := range themes {
		// Get theme name
		themeName := val["theme.Name."+baseLocale]
		if themeName == "" {
			return nil, nil, errors.New("web: empty theme.Name." + baseLocale + " parameter in '" + key + "' theme")
		}

		// Append to the translation, if it is not complete
		defTheme := themes[baseTheme]
		defTotal := len(defTheme)
		curTotal := 0
		for defKey, defVal := range defTheme {
			_, isExist := val[defKey]
			if isExist {
				curTotal = curTotal + 1
			} else {
				if strings.HasPrefix(defKey, "theme.Name.") {
					val[defKey] = val["theme.Name."+baseLocale]
				} else {
					val[defKey] = defVal
				}
			}
		}

		if curTotal == 0 {
			return nil, nil, errors.New("web: theme '" + key + "' is empty")
		}

		// Add theme to themes list
		themeNameSuffix := ""
		if curTotal != defTotal {
			themeNameSuffix = fmt.Sprintf(" (%.2f%%)", (float32(curTotal)/float32(defTotal))*100)
		}
		themesList[baseLocale][key] = themeName + themeNameSuffix

		for localeCode, _ := range localesList {
			result, ok := val["theme.Name."+localeCode]
			if ok {
				themesList[localeCode][key] = result + themeNameSuffix
			} else {
				themesList[localeCode][key] = themeName + themeNameSuffix
			}
		}
	}

	// Check default theme exist
	_, ok := themes[defaultTheme]
	if ok == false {
		return nil, nil, errors.New("web: default theme '" + defaultTheme + "' not found")
	}

	return themes, themesList, nil
}

func (themesList ThemesList) getForLocale(req *http.Request) ThemesListPart {
	// Get theme by cookie
	langCookie := getCookie(req, "lang")
	if langCookie != "" {
		theme, ok := themesList[langCookie]
		if ok == true {
			return theme
		}
	}

	// Load default part theme
	theme, _ := themesList[baseLocale]
	return theme
}

func (themes Themes) findTheme(req *http.Request, defaultTheme string) Theme {
	// Get theme by cookie
	themeCookie := getCookie(req, "theme")
	if themeCookie != "" {
		theme, ok := themes[themeCookie]
		if ok == true {
			return theme
		}
	}

	// Load default theme
	theme, _ := themes[defaultTheme]
	return theme
}

func (theme Theme) theme(s string) string {
	for key, val := range theme {
		if key == s {
			return val
		}
	}

	panic(errors.New("web: theme: unknown theme key: " + s))
}

func (theme Theme) tryHighlight(source string, lexer string) template.HTML {
	return tryHighlight(source, lexer, theme.theme("highlight.Theme"))
}
