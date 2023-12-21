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

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lcomrade/lenpaste/internal/apiv1"
	"github.com/lcomrade/lenpaste/internal/cli"
	"github.com/lcomrade/lenpaste/internal/config"
	"github.com/lcomrade/lenpaste/internal/logger"
	"github.com/lcomrade/lenpaste/internal/netshare"
	"github.com/lcomrade/lenpaste/internal/raw"
	"github.com/lcomrade/lenpaste/internal/storage"
	"github.com/lcomrade/lenpaste/internal/web"
)

var Version = "unknown"

func readFile(path string) (string, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file
	fileByte, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Return result
	return string(fileByte), nil
}

func exitOnError(e error) {
	fmt.Fprintln(os.Stderr, "error:", e.Error())
	os.Exit(1)
}

func main() {
	var err error

	// Read environment variables and CLI flags
	c := cli.New(Version)

	flagAddress := c.AddStringVar("address", ":80", "HTTP server ADDRESS:PORT.", nil)

	flagDbDriver := c.AddStringVar("db-driver", "sqlite3", "Currently supported drivers: \"sqlite3\" and \"postgres\".", nil)
	flagDbSource := c.AddStringVar("db-source", "", "DB source.", &cli.FlagOptions{Required: true})
	flagDbMaxOpenConns := c.AddIntVar("db-max-open-conns", 25, "Maximum number of connections to the database.", nil)
	flagDbMaxIdleConns := c.AddIntVar("db-max-idle-conns", 5, "Maximum number of idle connections to the database.", nil)
	flagDbCleanupPeriod := c.AddDurationVar("db-cleanup-period", "1m", "Interval at which the DB is cleared of expired but not yet deleted pastes.", nil)

	flagRobotsDisallow := c.AddBoolVar("robots-disallow", "Prohibits search engine crawlers from indexing site using robots.txt file.")

	flagTitleMaxLen := c.AddIntVar("title-max-length", 100, "Maximum length of the paste title. If 0 disable title, if -1 disable length limit.", nil)
	flagBodyMaxLen := c.AddIntVar("body-max-length", 20000, "Maximum length of the paste body. If -1 disable length limit. Can't be -1.", nil)
	flagMaxLifetime := c.AddDurationVar("max-paste-lifetime", "unlimited", "Maximum lifetime of the paste. Examples: 10m, 1h 30m, 12h, 1w, 30d, 365d.", &cli.FlagOptions{
		PreHook: func(s string) (string, error) {
			if s == "never" || s == "unlimited" {
				return "", nil
			}

			return s, nil
		},
	})

	flagGetPastesPer5Min := c.AddUintVar("get-pastes-per-5min", 50, "Maximum number of pastes that can be VIEWED in 5 minutes from one IP. If 0 disable rate-limit.", nil)
	flagGetPastesPer15Min := c.AddUintVar("get-pastes-per-15min", 100, "Maximum number of pastes that can be VIEWED in 15 minutes from one IP. If 0 disable rate-limit.", nil)
	flagGetPastesPer1Hour := c.AddUintVar("get-pastes-per-1hour", 500, "Maximum number of pastes that can be VIEWED in 1 hour from one IP. If 0 disable rate-limit.", nil)
	flagNewPastesPer5Min := c.AddUintVar("new-pastes-per-5min", 15, "Maximum number of pastes that can be CREATED in 5 minutes from one IP. If 0 disable rate-limit.", nil)
	flagNewPastesPer15Min := c.AddUintVar("new-pastes-per-15min", 30, "Maximum number of pastes that can be CREATED in 15 minutes from one IP. If 0 disable rate-limit.", nil)
	flagNewPastesPer1Hour := c.AddUintVar("new-pastes-per-1hour", 40, "Maximum number of pastes that can be CREATED in 1 hour from one IP. If 0 disable rate-limit.", nil)

	flagServerAbout := c.AddStringVar("server-about", "", "Path to the TXT file that contains the server description.", nil)
	flagServerRules := c.AddStringVar("server-rules", "", "Path to the TXT file that contains the server rules.", nil)
	flagServerTerms := c.AddStringVar("server-terms", "", "Path to the TXT file that contains the server terms of use.", nil)

	flagAdminName := c.AddStringVar("admin-name", "", "Name of the administrator of this server.", nil)
	flagAdminMail := c.AddStringVar("admin-mail", "", "Email of the administrator of this server.", nil)

	flagUiDefaultLifetime := c.AddStringVar("ui-default-lifetime", "", "Lifetime of paste will be set by default in WEB interface. Examples: 10min, 1h, 1d, 2w, 6mon, 1y.", nil)
	flagUiDefaultTheme := c.AddStringVar("ui-default-theme", "dark", "Sets the default theme for the WEB interface. Examples: dark, light, my_theme.", nil)
	flagUiThemesDir := c.AddStringVar("ui-themes-dir", "", "Loads external WEB interface themes from directory.", nil)

	flagLenPasswdFile := c.AddStringVar("lenpasswd-file", "", "File in LenPasswd format. If set, authorization will be required to create pastes.", nil)

	c.Parse()

	// -body-max-length flag
	if *flagBodyMaxLen == 0 {
		exitOnError(errors.New("maximum body length cannot be 0"))
	}

	// -max-paste-lifetime
	maxLifeTime := int64(-1)

	if *flagMaxLifetime != 0 && *flagMaxLifetime < 600 {
		exitOnError(errors.New("maximum paste lifetime flag cannot have a value less than 10 minutes"))
		maxLifeTime = int64(*flagMaxLifetime / time.Second)
	}

	// Load server about
	serverAbout := ""
	if *flagServerAbout != "" {
		serverAbout, err = readFile(*flagServerAbout)
		if err != nil {
			exitOnError(err)
		}
	}

	// Load server rules
	serverRules := ""
	if *flagServerRules != "" {
		serverRules, err = readFile(*flagServerRules)
		if err != nil {
			exitOnError(err)
		}
	}

	// Load server "terms of use"
	serverTermsOfUse := ""
	if *flagServerTerms != "" {
		if serverRules == "" {
			exitOnError(errors.New("in order to set the Terms of Use you must also specify the Server Rules"))
		}

		serverTermsOfUse, err = readFile(*flagServerTerms)
		if err != nil {
			exitOnError(err)
		}
	}

	// Settings
	log := logger.New("2006/01/02 15:04:05")

	db, err := storage.NewPool(*flagDbDriver, *flagDbSource, *flagDbMaxOpenConns, *flagDbMaxIdleConns)
	if err != nil {
		exitOnError(err)
	}

	cfg := config.Config{
		Log:               log,
		RateLimitGet:      netshare.NewRateLimitSystem(*flagGetPastesPer5Min, *flagGetPastesPer15Min, *flagGetPastesPer1Hour),
		RateLimitNew:      netshare.NewRateLimitSystem(*flagNewPastesPer5Min, *flagNewPastesPer15Min, *flagNewPastesPer1Hour),
		Version:           Version,
		TitleMaxLen:       *flagTitleMaxLen,
		BodyMaxLen:        *flagBodyMaxLen,
		MaxLifeTime:       maxLifeTime,
		ServerAbout:       serverAbout,
		ServerRules:       serverRules,
		ServerTermsOfUse:  serverTermsOfUse,
		AdminName:         *flagAdminName,
		AdminMail:         *flagAdminMail,
		RobotsDisallow:    *flagRobotsDisallow,
		UiDefaultLifetime: *flagUiDefaultLifetime,
		UiDefaultTheme:    *flagUiDefaultTheme,
		UiThemesDir:       *flagUiThemesDir,
		LenPasswdFile:     *flagLenPasswdFile,
	}

	apiv1Data := apiv1.Load(db, cfg)

	rawData := raw.Load(db, cfg)

	// Init data base
	err = storage.InitDB(*flagDbDriver, *flagDbSource)
	if err != nil {
		exitOnError(err)
	}

	// Load pages
	webData, err := web.Load(db, cfg)
	if err != nil {
		exitOnError(err)
	}

	// Handlers
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		webData.Handler(rw, req)
	})
	http.HandleFunc("/raw/", func(rw http.ResponseWriter, req *http.Request) {
		rawData.Hand(rw, req)
	})
	http.HandleFunc("/api/", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.Hand(rw, req)
	})

	// Run background job
	go func(cleanJobPeriod time.Duration) {
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
	}(*flagDbCleanupPeriod)

	// Run HTTP server
	log.Info("Run HTTP server on " + *flagAddress)
	err = http.ListenAndServe(*flagAddress, nil)
	if err != nil {
		exitOnError(err)
	}
}
