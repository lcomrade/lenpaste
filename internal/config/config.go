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
	"git.lcomrade.su/root/lenpaste/internal/model"
)

type Config struct {
	HTTP   ConfigHTTP   `json:"http"`
	DB     ConfigDB     `json:"database"`
	Public ConfigPublic `json:"public"`
	UI     ConfigUI     `json:"ui"`
	Paste  ConfigPaste  `json:"paste"`
	Auth   ConfigAuth   `json:"auth"`

	About      map[string]string `json:"-"`
	Rules      map[string]string `json:"-"`
	TermsOfUse map[string]string `json:"-"`

	Paths ConfigPaths `json:"-"`
}

type ConfigHTTP struct {
	Address string `json:"address"`
}

type ConfigDB struct {
	Driver             string `json:"driver"`
	Source             string `json:"source"`
	MaxOpenConns       int    `json:"max_open_conns"`
	MaxIdleConns       int    `json:"max_idle_conns"`
	ConnMaxLifetime    int64  `json:"-"`
	ConnMaxLifetimeStr string `json:"conn_max_lifetime"`

	CleanupPeriod    int64  `json:"-"`
	CleanupPeriodStr string `json:"cleanup_period"`
}

type ConfigPublic struct {
	AdminName string `json:"admin_name"`
	AdminMail string `json:"admin_mail"`

	RobotsDisallow bool `json:"robots_disallow"`
}

type ConfigUI struct {
	DefaultTheme string `json:"default_theme"`
}

type ConfigPaste struct {
	TitleMaxLen    int    `json:"title_max_len"`
	BodyMaxLen     int    `json:"body_max_len"`
	MaxLifetime    int64  `json:"-"`
	MaxLifetimeStr string `json:"max_lifetime"`

	UiDefaultLifetime string `json:"ui_default_lifetime"`
}

type ConfigAuth struct {
	Method string `json:"method"`
}

type ConfigPaths struct {
	MainCfg string `json:"-"`

	AboutDir string `json:"-"`
	RulesDir string `json:"-"`
	TermsDir string `json:"-"`

	ThemesDir string `json:"-"`

	LenPasswdFile string `json:"-"`
}

func (cfg *Config) GetAbout(locale string) string {
	out, ok := cfg.About[locale]
	if !ok {
		out, ok = cfg.About[model.BaseLocale]
		if !ok {
			return ""
		}
	}

	return out
}

func (cfg *Config) GetRules(locale string) string {
	out, ok := cfg.Rules[locale]
	if !ok {
		out, ok = cfg.Rules[model.BaseLocale]
		if !ok {
			return ""
		}
	}

	return out
}

func (cfg *Config) GetTermsOfUse(locale string) string {
	out, ok := cfg.TermsOfUse[locale]
	if !ok {
		out, ok = cfg.TermsOfUse[model.BaseLocale]
		if !ok {
			return ""
		}
	}

	return out
}
