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
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/apiv1"
	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/netshare"
	"git.lcomrade.su/root/lenpaste/internal/raw"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"git.lcomrade.su/root/lenpaste/internal/web"
)

func backgroundJob(cleanJobPeriod time.Duration, db storage.DB, log logger.Logger) {
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

func readFile(path string) (string, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file
	fileByte, err := ioutil.ReadAll(file)
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

func printHelp(noErrors bool) {
	println("Usage:", os.Args[0], "[-db-source] [OPTION]...")
	println("")
	println("  -address                ADDRESS:PORT (default: :80)")
	println()
	println("  -db-driver              Currently supported drivers: 'sqlite3' and 'postgres' (default: sqlite3)")
	println("  -db-source              DB source.")
	println("  -db-max-open-conns      Maximum number of connections to the database. (default: 25)")
	println("  -db-max-idle-conns      Maximum number of idle connections to the database. (default: 5)")
	println("  -db-cleanup-period      Interval at which the DB is cleared of expired but not yet deleted pastes. (default: 3h)")
	println()
	println("  -robots-disallow        Prohibits search engine crawlers from indexing site using robots.txt file.")
	println()
	println("  -title-max-length       Maximum length of the paste title. If 0 disable title, if -1 disable length limit. (default: 100)")
	println("  -body-max-length        Maximum length of the paste body. If -1 disable length limit. Can't be -1. (default: 20000)")
	println("  -max-paste-lifetime     Maximum lifetime of the paste. Examples: 10m, 1h 30m, 12h, 1w, 30d, 365d. (default: unlimited)")
	println()
	println("  -get-pastes-per-5min    Maximum number of pastes that can be VIEWED in 5 minutes from one IP. If 0 disable rate-limit. (default: 50)")
	println("  -get-pastes-per-15min   Maximum number of pastes that can be VIEWED in 15 minutes from one IP. If 0 disable rate-limit. (default: 100)")
	println("  -get-pastes-per-1hour   Maximum number of pastes that can be VIEWED in 1 hour from one IP. If 0 disable rate-limit. (default: 500)")
	println("  -new-pastes-per-5min    Maximum number of pastes that can be CREATED in 5 minutes from one IP. If 0 disable rate-limit. (default: 15)")
	println("  -new-pastes-per-15min   Maximum number of pastes that can be CREATED in 15 minutes from one IP. If 0 disable rate-limit. (default: 30)")
	println("  -new-pastes-per-1hour   Maximum number of pastes that can be CREATED in 1 hour from one IP. If 0 disable rate-limit. (default: 40)")
	println()
	println("  -server-about           Path to the TXT file that contains the server description.")
	println("  -server-rules           Path to the TXT file that contains the server rules.")
	println("  -server-terms           Path to the TXT file that contains the server terms of use.")
	println()
	println("  -admin-name             Name of the administrator of this server.")
	println("  -admin-mail             Email of the administrator of this server.")
	println()
	println("  -ui-default-lifetime    Lifetime of paste will be set by default in WEB interface. Examples: 10min, 1h, 1d, 2w, 6mon, 1y.")
	println("  -ui-default-theme       Sets the default theme for the WEB interface. Examples: dark, light, my_theme. (default: dark)")
	println("  -ui-themes-dir          Loads external WEB interface themes from directory.")
	println()
	println("  -lenpasswd-file         File in LenPasswd format. If set, authorization will be required to create pastes.")
	println()
	println("  -version                Display version and exit.")
	println("  -help                   Display this help and exit.")
	println()
	println("Exit status:")
	println(" 0  if you used the -help or -version flag")
	println(" 1  if there is an error during initialization")
	println(" 2  if command line arguments are not entered correctly")

	if noErrors == false {
		os.Exit(2)
	}

	os.Exit(0)
}

func printVersion() {
	println(model.Version)

	os.Exit(0)
}

func printFlagNotSet(flg string) {
	println("flag is not set:", flg)

	os.Exit(2)
}

func main() {
	var err error

	// Read cmd args
	flag.Usage = func() { printHelp(false) }

	flagAddress := flag.String("address", ":80", "")

	flagDbDriver := flag.String("db-driver", "sqlite3", "")
	flagDbSource := flag.String("db-source", "", "")
	flagDbMaxOpenConns := flag.Int("db-max-open-conns", 25, "")
	flagDbMaxIdleConns := flag.Int("db-max-idle-conns", 5, "")
	flagDbCleanupPeriod := flag.String("db-cleanup-period", "3h", "")

	flagRobotsDisallow := flag.Bool("robots-disallow", false, "")

	flagTitleMaxLen := flag.Int("title-max-length", 100, "")
	flagBodyMaxLen := flag.Int("body-max-length", 20000, "")
	flagMaxLifetime := flag.String("max-paste-lifetime", "unlimited", "")

	flagGetPastesPer5Min := flag.Uint("get-pastes-per-5min", 50, "")
	flagGetPastesPer15Min := flag.Uint("get-pastes-per-15min", 100, "")
	flagGetPastesPer1Hour := flag.Uint("get-pastes-per-1hour", 500, "")
	flagNewPastesPer5Min := flag.Uint("new-pastes-per-5min", 15, "")
	flagNewPastesPer15Min := flag.Uint("new-pastes-per-15min", 30, "")
	flagNewPastesPer1Hour := flag.Uint("new-pastes-per-1hour", 40, "")

	flagServerAbout := flag.String("server-about", "", "")
	flagServerRules := flag.String("server-rules", "", "")
	flagServerTerms := flag.String("server-terms", "", "")

	flagAdminName := flag.String("admin-name", "", "")
	flagAdminMail := flag.String("admin-mail", "", "")

	flagUiDefaultLifetime := flag.String("ui-default-lifetime", "", "")
	flagUiDefaultTheme := flag.String("ui-default-theme", "dark", "")
	flagUiThemesDir := flag.String("ui-themes-dir", "", "")

	flagLenPasswdFile := flag.String("lenpasswd-file", "", "")

	flagVersion := flag.Bool("version", false, "")
	flagHelp := flag.Bool("help", false, "")

	flag.Parse()

	// -help flag
	if *flagHelp == true {
		printHelp(true)
	}

	// -version flag
	if *flagVersion == true {
		printVersion()
	}

	// -db-source flag
	if *flagDbSource == "" {
		printFlagNotSet("-db-source")
	}

	// -body-max-length flag
	if *flagBodyMaxLen == 0 {
		println("-body-max-length flag cannot be 0")
		os.Exit(2)
	}

	// -max-paste-lifetime
	maxLifeTime := int64(-1)

	if *flagMaxLifetime != "never" && *flagMaxLifetime != "unlimited" {
		maxLifeTime, err = parseDuration(*flagMaxLifetime)
		if err != nil {
			exitOnError(err)
		}

		if maxLifeTime < 600 {
			println("-max-paste-lifetime flag cannot have a value less than 10 minutes")
			os.Exit(2)
		}
	}

	// -new-pastes-per-5min flag
	if *flagNewPastesPer5Min < 0 {
		println("-new-pastes-per-5min flag cannot be negative")
		os.Exit(2)
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
			println("In order to set the Terms of Use you must also specify the Server Rules. Server rules - this is a document written clearly for ordinary users. A Terms of Use is needed to protect the owner of the server from legal problems.")
			os.Exit(2)
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

	cfg := &config.Config{
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

	apiv1Data := apiv1.Load(log, db, cfg)

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
	jobDuration, err := time.ParseDuration(*flagDbCleanupPeriod)
	if err != nil {
		exitOnError(err)
	}

	go backgroundJob(jobDuration, db, log)

	// Run HTTP server
	log.Info("Run HTTP server on " + *flagAddress)
	err = http.ListenAndServe(*flagAddress, nil)
	if err != nil {
		exitOnError(err)
	}
}
