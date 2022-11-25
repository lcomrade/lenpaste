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
	"net/http"
	"path/filepath"
	"strings"
)

const defaultTheme = "dark"

type Theme map[string]string
type Themes map[string]Theme

func loadThemes(f embed.FS, themeDir string) (Themes, map[string]string, error) {
	themes := make(Themes)
	themesList := make(map[string]string)

	// Get theme files list
	files, err := f.ReadDir(themeDir)
	if err != nil {
		return nil, nil, errors.New("web: failed read dir '" + themeDir + "': " + err.Error())
	}

	// Load themes
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
		fileByte, err := f.ReadFile(filePath)
		if err != nil {
			return nil, nil, errors.New("web: failed open file '" + filePath + "': " + err.Error())
		}

		fileStr := bytes.NewBuffer(fileByte).String()

		// Load theme
		theme, err := readKVCfg(fileStr)
		if err != nil {
			return nil, nil, errors.New("web: failed read file '" + filePath + "': " + err.Error())
		}

		themes[themeCode] = Theme(theme)
	}

	// Prepare themes list
	for key, val := range themes {
		// Get theme name
		themeName := val["theme.Name"]
		if themeName == "" {
			return nil, nil, errors.New("web: empty theme.Name parameter in '" + key + "' theme")
		}

		// Append to the translation, if it is not complete
		defTheme := themes[defaultTheme]
		defTotal := len(defTheme)
		curTotal := 0
		for defKey, defVal := range defTheme {
			_, isExist := val[defKey]
			if isExist {
				curTotal = curTotal + 1
			} else {
				val[defKey] = defVal
			}
		}

		if curTotal == 0 {
			return nil, nil, errors.New("web: theme '" + key + "' is empty")
		}

		themesList[key] = themeName + fmt.Sprintf(" (%.2f%%)", (float32(curTotal)/float32(defTotal))*100)
	}

	return themes, themesList, nil
}

func (themes Themes) findTheme(req *http.Request) Theme {
	// Get accept language by cookie
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

	return s
}
