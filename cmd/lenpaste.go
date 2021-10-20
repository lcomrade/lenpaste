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
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	//Background job break
	backJobBreak = 1 * time.Minute
	//Logs files
	logDir      = "./data/log"
	logFileMod  = 0700
	logInfoFile = "./data/log/info"
	logErrFile  = "./data/log/errors"
	logJobFile  = "./data/log/job"
)

func tee(file string, writer io.Writer) (io.Writer, error) {
	var mw io.Writer

	//Open file
	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, logFileMod)
	if err != nil {
		return mw, err
	}
	//defer logFile.Close()

	//Create Multi Writer
	mw = io.MultiWriter(writer, logFile)

	//Return
	return mw, nil
}

func BackgroundJob() {
	//Prepare loging
	logMW, err := tee(logJobFile, os.Stderr)
	if err != nil {
		panic(err)
	}

	backLog := log.New(logMW, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Run
	for {
		//Delete expired pastes
		errs := storage.DelExpiredPaste()
		if len(errs) != 0 {
			for _, err := range errs {
				backLog.Println("delete expired error:", err)
			}
		}

		//Wait
		time.Sleep(backJobBreak)

	}
}

func init() {
	//Create log dir
	err := os.MkdirAll(logDir, logFileMod)
	if err != nil {
		panic(err)
	}
}

func main() {
	//Prepare loging (ERROR)
	logErrMW, err := tee(logErrFile, os.Stderr)
	if err != nil {
		panic(err)
	}

	errLog := log.New(logErrMW, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Prepare loging (INFO)
	logInfoMW, err := tee(logInfoFile, os.Stdout)
	if err != nil {
		panic(err)
	}
	infoLog := log.New(logInfoMW, "INFO\t", log.Ldate|log.Ltime)

	//Log start
	infoLog.Println("## Start Lenpaste ##")

	//Read config
	config, err := config.ReadConfig()
	if err != nil {
		errLog.Fatal(err)
	}

	//Load assets
	err = assets.Load()
	if err != nil {
		errLog.Fatal(err)
	}

	//Load pages
	err = pages.Load()
	if err != nil {
		errLog.Fatal(err)
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
	http.HandleFunc("/api/about", api.GetAbout)
	http.HandleFunc("/api/rules", api.GetRules)
	http.HandleFunc("/api/version", api.GetVersion)

	//Run (Background Job)
	go BackgroundJob()

	//Run (WEB)
	infoLog.Println("HTTP server listen: '" + config.HTTP.Listen + "'")
	infoLog.Println("Use TLS: ", config.HTTP.UseTLS)

	//Use TLS or no?
	if config.HTTP.UseTLS == false {
		//No TLS
		err = http.ListenAndServe(config.HTTP.Listen, nil)
		if err != nil {
			errLog.Fatal(err)
		}

	} else {
		//Enable TLS
		infoLog.Println("SSL cert: '" + config.HTTP.SSLCert + "'")
		infoLog.Println("SSL key: '" + config.HTTP.SSLKey + "'")

		err = http.ListenAndServeTLS(config.HTTP.Listen, config.HTTP.SSLCert, config.HTTP.SSLKey, nil)
		if err != nil {
			errLog.Fatal(err)
		}
	}
}
