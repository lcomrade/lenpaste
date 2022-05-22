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
	"git.lcomrade.su/root/lenpaste/internal/apiv1"
	"git.lcomrade.su/root/lenpaste/internal/web"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"errors"
	"net/http"
	"time"
	"flag"
	"os"
)

func backgroundJob(cleanJobPeriod time.Duration, db storage.DB, log logger.Config) {
	for {
		// Delete expired pastes
		count, err := db.PasteDeleteExpired()
		if err != nil {
			log.Error(errors.New("Delete expired: " + err.Error()))
		}

		log.Info("Deleted "+string(count)+" expired pastes")

		// Wait
		time.Sleep(cleanJobPeriod)
	}
}


func printHelp() {
	println("Usage:", os.Args[0], "[OPTION]...")
	println("")
	println("    --db-source path to config file")
	println("-h, --help      display this help and exit")

	os.Exit(0)
}

func printFlagNotSet(flg string) {
	println("flag is not set:", flg)

	os.Exit(0)
}

func main() {
	// Read cmd args
	flag.Usage = printHelp

	flagDbSource := flag.String("db-source", "", "")
	flagH := flag.Bool("h", false, "")
	flagHelp := flag.Bool("-help", false, "")

	flag.Parse()

	// -h or --help flag
	if *flagH == true || *flagHelp == true {
		printHelp()
	}

	// --db-source flag
	if *flagDbSource == "" {
		printFlagNotSet("--db-source")
	}

	// Settings
	db := storage.DB {
		DriverName: "sqlite3",
		DataSourceName: *flagDbSource,
	}

	log := logger.Config{
		TimeFormat: "2006/01/02 15:04:05",
	}

	apiv1Data := apiv1.Data{
		DB: db,
		Log: log,
	}

	// Init data base
	err := db.InitDB()
	if err != nil {
		panic(err)
	}

	// Load pages
	webData, err := web.Load("./web", db, log)
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
	http.HandleFunc("/new", func(rw http.ResponseWriter, req *http.Request) {
		webData.NewHand(rw, req)
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
	log.Info("Run background job")
	go backgroundJob(10 * time.Minute, db, log)

	// Run HTTP server
	log.Info("Run HTTP server on :8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
