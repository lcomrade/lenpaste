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

package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

func Load(cfgDir string) (*Config, error) {
	// Default configuration
	cfg := Config{
		HTTP: ConfigHTTP{
			Address: ":80",
		},

		DB: ConfigDB{
			Driver:             "",
			Source:             "",
			MaxOpenConns:       25,
			MaxIdleConns:       5,
			ConnMaxLifetime:    5 * 60,
			ConnMaxLifetimeStr: "5m",

			CleanupPeriod:    60 * 60 * 3,
			CleanupPeriodStr: "3h",
		},

		Public: ConfigPublic{
			AdminName: "",
			AdminMail: "",

			RobotsDisallow: false,
		},

		UI: ConfigUI{
			DefaultTheme: "dark",
		},

		Paste: ConfigPaste{
			TitleMaxLen:    100,
			BodyMaxLen:     20000,
			MaxLifetime:    0,
			MaxLifetimeStr: "",

			UiDefaultLifetime: "",
		},

		Auth: ConfigAuth{
			Method: "",
		},

		About:      nil,
		Rules:      nil,
		TermsOfUse: nil,

		Paths: ConfigPaths{
			MainCfg: filepath.Join(cfgDir, model.SmallName+".json"),

			AboutDir: filepath.Join(cfgDir, "about"),
			RulesDir: filepath.Join(cfgDir, "rules"),
			TermsDir: filepath.Join(cfgDir, "terms"),

			ThemesDir: filepath.Join(cfgDir, "themes"),

			LenPasswdFile: filepath.Join(cfgDir, "lenpasswd"),
		},
	}

	// Read main configuration file
	cfgFile, err := os.Open(cfg.Paths.MainCfg)
	if err != nil {
		return nil, errors.New("config: " + err.Error())
	}
	defer cfgFile.Close()

	err = json.NewDecoder(cfgFile).Decode(&cfg)
	if err != nil {
		return nil, errors.New("config: " + err.Error())
	}

	// Convert strings duration to time
	cfg.DB.ConnMaxLifetime, err = parseDuration(cfg.DB.ConnMaxLifetimeStr)
	if err != nil {
		return nil, errors.New("config: " + err.Error())
	}

	cfg.Paste.MaxLifetime, err = parseDuration(cfg.Paste.MaxLifetimeStr)
	if err != nil {
		return nil, errors.New("config: " + err.Error())
	}

	// Read about, rules and terms of use
	cfg.About, err = loadL10nFiles(cfg.Paths.AboutDir, ".txt")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("config: " + err.Error())
		}
	}

	cfg.Rules, err = loadL10nFiles(cfg.Paths.RulesDir, ".txt")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("config: " + err.Error())
		}
	}

	cfg.TermsOfUse, err = loadL10nFiles(cfg.Paths.TermsDir, ".txt")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("config: " + err.Error())
		}
	}

	return &cfg, nil
}

func loadL10nFiles(dir, ext string) (map[string]string, error) {
	// Get list files in directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Read files
	out := make(map[string]string)

	for _, part := range files {
		if part.IsDir() {
			continue
		}

		fileName := part.Name()
		if !strings.HasSuffix(fileName, ext) {
			continue
		}

		// Read file and add it to map
		fileByte, err := os.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		out[strings.TrimSuffix(fileName, ext)] = bytes.NewBuffer(fileByte).String()
	}

	// If map is empty return nil
	if len(out) == 0 {
		return nil, nil
	}

	return out, nil
}
