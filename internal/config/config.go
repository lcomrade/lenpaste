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
	"github.com/lcomrade/lenpaste/internal/logger"
	"github.com/lcomrade/lenpaste/internal/netshare"
)

const Software = "Lenpaste"

type Config struct {
	Log logger.Logger

	RateLimitNew *netshare.RateLimitSystem
	RateLimitGet *netshare.RateLimitSystem

	Version string

	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsOfUse string

	AdminName string
	AdminMail string

	RobotsDisallow bool

	LenPasswdFile string

	UiDefaultLifetime string
	UiDefaultTheme    string
	UiThemesDir       string
}
