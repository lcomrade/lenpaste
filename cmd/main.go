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
	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/raw"
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"git.lcomrade.su/root/lenpaste/internal/web"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var Version = ""

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
	println("error:", e.Error())
	os.Exit(1)
}

func printHelp(noErrors bool) {
	println("Usage:", os.Args[0], "[-web-dir] [OPTION]...")
	println("")
	println("  -addres             ADDRES:PORT (default: :80)")
	println("  -web-dir            Dir with page templates and static content")
	println("  -db-driver          Currently supported drivers: 'sqlite3' and 'postgres' (default: sqlite3)")
	println("  -db-source          DB source")
	println("  -db-cleanup-period  Interval at which the DB is cleared of expired but not yet deleted pastes. (default: 3h)")
	println("  -robots-disallow    Prohibits search engine crawlers from indexing site using robots.txt file.")
	println("  -title-max-length   Maximum length of the paste title. If 0 disable title, if -1 disable length limit. (default: 100)")
	println("  -body-max-length    Maximum length of the paste body. If -1 disable length limit. Can't be -1. (default: 100000)")
	println("  -max-paste-lifetime Maximum lifetime of the paste. Examples: 12h, 7m, 10s. (default: never)")
	println("  -server-about       Path to the HTML file that contains the server description.")
	println("  -server-rules       Path to the HTML file that contains the server rules.")
	println("  -admin-name         Name of the administrator of this server")
	println("  -admin-mail         Email of the administrator of this server.")
	println("  -version            Display version and exit")
	println("  -help               Display this help and exit")
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
	println(Version)

	os.Exit(0)
}

func printFlagNotSet(flg string) {
	println("flag is not set:", flg)

	os.Exit(2)
}

func init() {
	if Version == "" {
		Version = "unknown"
	}
}

func main() {
	// Get ./bin/ dir
	binFile, err := os.Executable()
	if err != nil {
		exitOnError(err)
	}

	binFile, err = filepath.EvalSymlinks(binFile)
	if err != nil {
		exitOnError(err)
	}

	binDir := filepath.Dir(binFile)

	// Get ./share/lenpaste/web dir
	defaultWebDir := filepath.Join(binDir, "../share/lenpaste/web")

	// Read cmd args
	flag.Usage = func() { printHelp(false) }

	flagAddress := flag.String("address", ":80", "")
	flagWebDir := flag.String("web-dir", defaultWebDir, "")
	flagDbDriver := flag.String("db-driver", "sqlite3", "")
	flagDbSource := flag.String("db-source", "", "")
	flagDbCleanupPeriod := flag.String("db-cleanup-period", "3h", "")
	flagRobotsDisallow := flag.Bool("robots-disallow", false, "")
	flagTitleMaxLen := flag.Int("title-max-length", 100, "")
	flagBodyMaxLen := flag.Int("body-max-length", 10000, "")
	flagMaxLifetime := flag.String("max-paste-lifetime", "never", "")
	flagServerAbout := flag.String("server-about", "", "")
	flagServerRules := flag.String("server-rules", "", "")
	flagAdminName := flag.String("admin-name", "", "")
	flagAdminMail := flag.String("admin-mail", "", "")
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

	if *flagMaxLifetime != "never" {
		maxLifeTimeTmp, err := time.ParseDuration(*flagMaxLifetime)
		if err != nil {
			exitOnError(err)
		}

		maxLifeTime = int64(maxLifeTimeTmp.Seconds())

		if maxLifeTime < 600 {
			println("-max-paste-lifetime flag cannot have a value less than 10 minutes")
			os.Exit(2)
		}
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

	// Settings
	db := storage.DB{
		DriverName:     *flagDbDriver,
		DataSourceName: *flagDbSource,
	}

	log := logger.Config{
		TimeFormat: "2006/01/02 15:04:05",
	}

	cfg := config.Config{
		DB:             db,
		Log:            log,
		Version:        Version,
		TitleMaxLen:    *flagTitleMaxLen,
		BodyMaxLen:     *flagBodyMaxLen,
		MaxLifeTime:    maxLifeTime,
		ServerAbout:    serverAbout,
		ServerRules:    serverRules,
		AdminName:      *flagAdminName,
		AdminMail:      *flagAdminMail,
		RobotsDisallow: *flagRobotsDisallow,
	}

	apiv1Data := apiv1.Load(cfg)

	rawData := raw.Load(cfg)

	// Init data base
	err = db.InitDB()
	if err != nil {
		exitOnError(err)
	}

	// Load pages
	webData, err := web.Load(cfg, *flagWebDir)
	if err != nil {
		exitOnError(err)
	}

	// Handlers
	http.HandleFunc("/robots.txt", func(rw http.ResponseWriter, req *http.Request) {
		webData.RobotsTxtHand(rw, req)
	})
	if *flagRobotsDisallow == false {
		http.HandleFunc("/sitemap.xml", func(rw http.ResponseWriter, req *http.Request) {
			webData.SitemapHand(rw, req)
		})
	}

	http.HandleFunc("/style.css", func(rw http.ResponseWriter, req *http.Request) {
		webData.StyleCSSHand(rw, req)
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		webData.MainHand(rw, req)
	})
	http.HandleFunc("/main.js", func(rw http.ResponseWriter, req *http.Request) {
		webData.MainJSHand(rw, req)
	})
	http.HandleFunc("/paste.js", func(rw http.ResponseWriter, req *http.Request) {
		webData.PasteJSHand(rw, req)
	})

	http.HandleFunc("/raw/", func(rw http.ResponseWriter, req *http.Request) {
		rawData.RawHand(rw, req)
	})
	http.HandleFunc("/dl/", func(rw http.ResponseWriter, req *http.Request) {
		webData.DlHand(rw, req)
	})
	http.HandleFunc("/emb/", func(rw http.ResponseWriter, req *http.Request) {
		webData.EmbeddedHand(rw, req)
	})
	http.HandleFunc("/emb_help/", func(rw http.ResponseWriter, req *http.Request) {
		webData.EmbeddedHelpHand(rw, req)
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
	http.HandleFunc("/docs/api_libs", func(rw http.ResponseWriter, req *http.Request) {
		webData.DocsApiLibsHand(rw, req)
	})

	http.HandleFunc("/api/", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.MainHand(rw, req)
	})
	http.HandleFunc("/api/v1/new", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.NewHand(rw, req)
	})
	http.HandleFunc("/api/v1/get", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.GetHand(rw, req)
	})
	http.HandleFunc("/api/v1/getServerInfo", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.GetServerInfoHand(rw, req)
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
