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
	"../internal/assets"
	"../internal/config"
	"../internal/pages"
	"../internal/storage"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	//Background job break
	backJobBreak = 1 * time.Minute
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
	//Read config
	config, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	//Load assets
	err = assets.Load()
	if err != nil {
		panic(err)
	}

	//Load pages
	err = pages.Load()
	if err != nil {
		panic(err)
	}

	//Handlers
	http.HandleFunc("/style.css", pages.Style)
	http.HandleFunc("/robots.txt", assets.RobotsTxt)

	http.HandleFunc("/", pages.GetPaste)
	http.HandleFunc("/new", pages.NewPaste)
	http.HandleFunc("/new_done", pages.NewPasteDone)
	http.HandleFunc("/api", pages.API)
	http.HandleFunc("/rules", pages.Rules)
	http.HandleFunc("/version", pages.Version)

	http.HandleFunc("/api/new", api.NewPaste)
	http.HandleFunc("/api/get/", api.GetPaste)
	http.HandleFunc("/api/rules", api.GetRules)

	//Run (Background Job)
	go BackgroundJob()

	//Run (WEB)
	logInfo.Println("HTTP server listen: '" + config.HTTP.Listen + "'")
	logInfo.Println("Use TLS: ", config.HTTP.UseTLS)

	//Use TLS or no?
	if config.HTTP.UseTLS == false {
		//No TLS
		err = http.ListenAndServe(config.HTTP.Listen, nil)
		if err != nil {
			panic(err)
		}

	} else {
		//Enable TLS
		logInfo.Println("SSL cert: '" + config.HTTP.SSLCert + "'")
		logInfo.Println("SSL key: '" + config.HTTP.SSLKey + "'")

		err = http.ListenAndServeTLS(config.HTTP.Listen, config.HTTP.SSLCert, config.HTTP.SSLKey, nil)
		if err != nil {
			panic(err)
		}
	}
}
