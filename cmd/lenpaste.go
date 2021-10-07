/*
   Copyright 2021 Leonid Maslakov

   License: GPL-3.0-or-later

   This file is part of Lenpaste.

   Lenpaste is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Lenpaste is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Lenpaste.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"../internal/api"
	"../internal/pages"
	"../internal/storage"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	//Background job break
	backJobBreak     = 1 * time.Minute
	httpServerListen = ":8000"
)

//Logging
var logInfo = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var logError = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func BackgroundJob() {
	for {
		//Delete expired pastes
		err := storage.DelExpiredPaste()
		if err != nil {
			logError.Println("Delete expired:", err)
		}

		//Wait
		time.Sleep(backJobBreak)

	}
}

func main() {
	//Pages
	err := pages.Load()
	if err != nil {
		panic(err)
	}

	//Handlers
	http.HandleFunc("/style.css", pages.Style)
	http.HandleFunc("/", pages.GetPaste)
	http.HandleFunc("/new", pages.NewPaste)
	http.HandleFunc("/new_done", pages.NewPasteDone)
	http.HandleFunc("/api/new", api.NewPaste)
	http.HandleFunc("/api/get/", api.GetPaste)

	//Run (Background Job)
	go BackgroundJob()

	//Run (WEB)
	logInfo.Println("HTTP server listen: '" + httpServerListen + "'")
	err = http.ListenAndServe(httpServerListen, nil)
	if err != nil {
		panic(err)
	}
}
