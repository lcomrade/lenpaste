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

package logger

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/netshare"
)

type Logger struct {
	timeFormat string
}

func New(timeFormat string) *Logger {
	return &Logger{
		timeFormat: timeFormat,
	}
}

func getTrace() string {
	trace := ""

	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			trace = trace + file + "#" + strconv.Itoa(line) + ": "

		} else {
			return trace
		}
	}
}

func (cfg *Logger) Info(msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format(cfg.timeFormat), "[INFO]   ", msg)
}

func (cfg *Logger) Error(e error) {
	fmt.Fprintln(os.Stderr, time.Now().Format(cfg.timeFormat), "[ERROR]  ", getTrace(), e.Error())
}

func (cfg *Logger) HttpRequest(req *http.Request, code int) {
	fmt.Fprintln(os.Stdout, time.Now().Format(cfg.timeFormat), "[REQUEST]", netshare.GetClientAddr(req).String(), req.Method, code, req.URL.Path, "(User-Agent: "+req.UserAgent()+")")
}

func (cfg *Logger) HttpError(req *http.Request, e error) {
	fmt.Fprintln(os.Stderr, time.Now().Format(cfg.timeFormat), "[ERROR]  ", netshare.GetClientAddr(req).String(), req.Method, 500, req.URL.Path, "(User-Agent: "+req.UserAgent()+")", "Error:", getTrace(), e.Error())
}
