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

package apiv1

import (
	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	chromaLexers "github.com/alecthomas/chroma/lexers"
)

type Data struct {
	Log logger.Config
	DB  storage.DB

	Lexers *[]string

	Version *string

	TitleMaxLen *int
	BodyMaxLen  *int
	MaxLifeTime *int64

	ServerAbout      *string
	ServerRules      *string
	ServerTermsOfUse *string

	AdminName *string
	AdminMail *string

	UiDefaultLifeTime *string
}

func Load(cfg config.Config) Data {
	lexers := chromaLexers.Names(false)

	return Data{
		DB:                cfg.DB,
		Log:               cfg.Log,
		Lexers:            &lexers,
		Version:           &cfg.Version,
		TitleMaxLen:       &cfg.TitleMaxLen,
		BodyMaxLen:        &cfg.BodyMaxLen,
		MaxLifeTime:       &cfg.MaxLifeTime,
		ServerAbout:       &cfg.ServerAbout,
		ServerRules:       &cfg.ServerRules,
		ServerTermsOfUse:  &cfg.ServerTermsOfUse,
		AdminName:         &cfg.AdminName,
		AdminMail:         &cfg.AdminMail,
		UiDefaultLifeTime: &cfg.UiDefaultLifetime,
	}
}
