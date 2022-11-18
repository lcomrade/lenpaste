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

package logger

import (
	"fmt"
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	TimeFormat string
}

func New(timeFormat string) Logger {
	return Logger{
		TimeFormat: timeFormat,
	}
}

func (cfg Logger) Info(msg string) {
	now := time.Now()
	fmt.Println(now.Format(cfg.TimeFormat), "[INFO]", msg)
}

func (cfg Logger) Error(err error) {
	now := time.Now()
	fmt.Fprintln(os.Stderr, now.Format(cfg.TimeFormat), "[ERROR]", err.Error())
}

func (cfg Logger) HttpRequest(req *http.Request) {
	now := time.Now()
	fmt.Println(now.Format(cfg.TimeFormat), "[REQUEST]", netshare.GetClientAddr(req), req.Method, req.URL.Path)
}

func (cfg Logger) HttpError(req *http.Request, err error) {
	now := time.Now()
	fmt.Fprintln(os.Stderr, now.Format(cfg.TimeFormat), "[ERROR]", netshare.GetClientAddr(req), req.Method, req.URL.Path, "Error:", err.Error())
}
