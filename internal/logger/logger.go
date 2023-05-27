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
	"os"
	"runtime"
	"strconv"
	"time"
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

func (log *Logger) Error(a ...any) {
	fmt.Fprint(os.Stderr, time.Now().Format(log.timeFormat), "     [ERROR]    ", getTrace(), fmt.Sprintln(a...))
}

func (log *Logger) Info(a ...any) {
	fmt.Fprint(os.Stdout, time.Now().Format(log.timeFormat), "     [INFO]     ", fmt.Sprintln(a...))
}

func (log *Logger) Warning(a ...any) {
	fmt.Fprint(os.Stdout, time.Now().Format(log.timeFormat), "     [WARNING]  ", fmt.Sprintln(a...))
}

func (log *Logger) Debug(a ...any) {
	fmt.Fprint(os.Stdout, time.Now().Format(log.timeFormat), "     [DEBUG]    ", fmt.Sprintln(a...))
}
