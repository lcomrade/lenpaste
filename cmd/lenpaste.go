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

package main

import (
	"git.lcomrade.su/root/lenpaste/internal/api"
	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/pages"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	//Logs files
	logDir       = "./data/log"
	logFileMod   = 0700
	logOldPrefix = ".old"
	logInfoFile  = "./data/log/info"
	logErrFile   = "./data/log/error"
	logJobFile   = "./data/log/job"
)

func tee(file string, writer io.Writer, save bool) (io.Writer, error) {
	//Check SAVE
	if save == false {
		return writer, nil
	}

	//If save == true: continue
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

func BackgroundJob(cleanJobPeriod time.Duration, saveLog bool) {
	//Prepare loging
	logWr, err := tee(logJobFile, os.Stderr, saveLog)
	if err != nil {
		panic(err)
	}

	backLog := log.New(logWr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
		time.Sleep(cleanJobPeriod)

	}
}

func rotateLog(path string, maxSize int64) error {
	//Stat file
	file, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) == true {
			return nil
		} else {
			return err
		}
	}

	//Check file size
	if file.Size() >= maxSize {
		//Move
		err := os.Rename(path, path+logOldPrefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	//Create log dir
	err := os.MkdirAll(logDir, logFileMod)
	if err != nil {
		panic(err)
	}
}

func main() {
	//Read config
	config, err := config.ReadConfig()
	if err != nil {
		panic(errors.New("read config: " + err.Error()))
	}

	if config.Logs.RotateLogs == true {
		//Rotate ERROR
		if config.Logs.SaveErr == true {
			err := rotateLog(logErrFile, config.Logs.MaxLogSize)
			if err != nil {
				panic(errors.New("rotate log: " + err.Error()))
			}
		}

		//Rotate INFO
		if config.Logs.SaveInfo == true {
			err := rotateLog(logInfoFile, config.Logs.MaxLogSize)
			if err != nil {
				panic(errors.New("rotate log: " + err.Error()))
			}
		}

		//Rotate JOB
		if config.Logs.SaveJob == true {
			err := rotateLog(logJobFile, config.Logs.MaxLogSize)
			if err != nil {
				panic(errors.New("rotate log: " + err.Error()))
			}
		}
	}

	//Prepare loging (ERROR)
	logErrWr, err := tee(logErrFile, os.Stderr, config.Logs.SaveErr)
	if err != nil {
		panic(err)
	}

	errLog := log.New(logErrWr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Prepare loging (INFO)
	logInfoWr, err := tee(logInfoFile, os.Stdout, config.Logs.SaveInfo)
	if err != nil {
		panic(err)
	}
	infoLog := log.New(logInfoWr, "INFO\t", log.Ldate|log.Ltime)

	//Log start
	infoLog.Println("## Start Lenpaste ##")

	//Load pages
	err = pages.Load()
	if err != nil {
		errLog.Fatal("load pages:", err)
	}

	//Handlers
	http.HandleFunc("/style.css", pages.Style)

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
	if config.Storage.EnableCleanJob == true {
		infoLog.Println("Run background clean job")
		go BackgroundJob(config.Storage.CleanJobPeriod, config.Logs.SaveJob)
	}

	//Run (WEB)
	infoLog.Println("HTTP server listen: '" + config.HTTP.Listen + "'")
	infoLog.Println("Use TLS: ", config.HTTP.UseTLS)

	//Use TLS or no?
	if config.HTTP.UseTLS == false {
		//No TLS
		err = http.ListenAndServe(config.HTTP.Listen, nil)
		if err != nil {
			errLog.Fatal("run WEB server:", err)
		}

	} else {
		//Enable TLS
		infoLog.Println("SSL cert: '" + config.HTTP.SSLCert + "'")
		infoLog.Println("SSL key: '" + config.HTTP.SSLKey + "'")

		err = http.ListenAndServeTLS(config.HTTP.Listen, config.HTTP.SSLCert, config.HTTP.SSLKey, nil)
		if err != nil {
			errLog.Fatal("run WEB server:", err)
		}
	}
}
