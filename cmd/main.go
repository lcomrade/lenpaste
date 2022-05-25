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
	"errors"
	"flag"
	"git.lcomrade.su/root/lenpaste/internal/apiv1"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/raw"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"git.lcomrade.su/root/lenpaste/internal/web"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func backgroundJob(cleanJobPeriod time.Duration, db storage.DB, log logger.Config) {
	for {
		// Delete expired pastes
		count, err := db.PasteDeleteExpired()
		if err != nil {
			log.Error(errors.New("Delete expired: " + err.Error()))
		}

		log.Info("Delete " + strconv.FormatInt(count, 10) + " expired pastes")

		// Wait
		time.Sleep(cleanJobPeriod)
	}
}

func printHelp() {
	println("Usage:", os.Args[0], "[OPTION]...")
	println("")
	println("-addres    ADDRES:PORT (default: :80)")
	println("-web-dir   Dir with page templates and static content")
	println("-db-driver Only 'sqlite3' is available yet (default: sqlite3)")
	println("-db-source DB source")
	println("-help      display this help and exit")

	os.Exit(0)
}

func printFlagNotSet(flg string) {
	println("flag is not set:", flg)

	os.Exit(0)
}

func main() {
	// Get ./bin/ dir
	binFile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	binFile, err = filepath.EvalSymlinks(binFile)
	if err != nil {
		panic(err)
	}

	binDir := filepath.Dir(binFile)

	// Get ./share/lenpaste/web dir
	defaultWebDir := filepath.Join(binDir, "../share/lenpaste/web")

	// Read cmd args
	flag.Usage = printHelp

	flagAddress := flag.String("address", ":80", "")
	flagWebDir := flag.String("web-dir", defaultWebDir, "")
	flagDbDriver := flag.String("db-driver", "sqlite3", "")
	flagDbSource := flag.String("db-source", "", "")
	flagHelp := flag.Bool("help", false, "")

	flag.Parse()

	// -help flag
	if *flagHelp == true {
		printHelp()
	}

	// -db-source flag
	if *flagDbSource == "" {
		printFlagNotSet("-db-source")
	}

	// Settings
	db := storage.DB{
		DriverName:     *flagDbDriver,
		DataSourceName: *flagDbSource,
	}

	log := logger.Config{
		TimeFormat: "2006/01/02 15:04:05",
	}

	apiv1Data := apiv1.Data{
		DB:  db,
		Log: log,
	}

	rawData := raw.Data{
		DB:  db,
		Log: log,
	}

	// Init data base
	err = db.InitDB()
	if err != nil {
		panic(err)
	}

	// Load pages
	webData, err := web.Load(*flagWebDir, db, log)
	if err != nil {
		panic(err)
	}

	// Handlers
	http.HandleFunc("/style.css", func(rw http.ResponseWriter, req *http.Request) {
		webData.StyleCSSHand(rw, req)
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		webData.MainHand(rw, req)
	})

	http.HandleFunc("/raw/", func(rw http.ResponseWriter, req *http.Request) {
		rawData.MainHand(rw, req)
	})

	http.HandleFunc("/about", func(rw http.ResponseWriter, req *http.Request) {
		webData.AboutHand(rw, req)
	})
	http.HandleFunc("/about/license", func(rw http.ResponseWriter, req *http.Request) {
		webData.LicenseHand(rw, req)
	})
	http.HandleFunc("/about/source_code", func(rw http.ResponseWriter, req *http.Request) {
		webData.SourceCodePageHand(rw, req)
	})

	http.HandleFunc("/docs", func(rw http.ResponseWriter, req *http.Request) {
		webData.DocsHand(rw, req)
	})
	http.HandleFunc("/docs/apiv1", func(rw http.ResponseWriter, req *http.Request) {
		webData.DocsApiV1Hand(rw, req)
	})

	http.HandleFunc("/api/v1/new", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.NewHand(rw, req)
	})
	http.HandleFunc("/api/v1/get", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.GetHand(rw, req)
	})

	// Run background job
	go backgroundJob(10*time.Minute, db, log)

	// Run HTTP server
	log.Info("Run HTTP server on " + *flagAddress)
	err = http.ListenAndServe(*flagAddress, nil)
	if err != nil {
		panic(err)
	}
}
