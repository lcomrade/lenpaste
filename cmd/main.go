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
	"os"
	"strconv"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/handler"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
	"git.lcomrade.su/root/lenpaste/internal/storage"

	"github.com/urfave/cli/v2"
)

func fatal(e error) {
	fmt.Fprintln(os.Stderr, model.SmallName+" error:", e.Error())
	os.Exit(1)
}

func backgroundJobPastes(cleanJobPeriod time.Duration, db *storage.DB, log *logger.Logger) {
	for {
		// Delete expired pastes
		count, err := db.PasteDeleteExpired()
		if err != nil {
			log.Error(errors.New("background: " + err.Error()))
		}

		log.Info("Delete " + strconv.FormatInt(count, 10) + " expired pastes.")

		// Wait
		time.Sleep(cleanJobPeriod)
	}
}

func backgroundJobFiles(cleanJobPeriod time.Duration, db *storage.DB, log *logger.Logger) {
	for {
		expired, notFinished, err := db.FileCleanup()
		if err != nil {
			log.Error(errors.New("background: " + err.Error()))
		}

		log.Info("Delete " + strconv.FormatInt(expired, 10) + " expired files.")
		log.Info("Delete " + strconv.FormatInt(notFinished, 10) + " unfinished uploads.")

		// Wait
		time.Sleep(cleanJobPeriod)
	}
}

func run(cfgDir string) error {
	// Setup logger
	log := logger.New("2006/01/02 15:04:05")

	// Read configurations files
	cfg, err := config.Load(cfgDir)
	if err != nil {
		return err
	}

	// Open and init database
	db, err := storage.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.InitDB()
	if err != nil {
		return err
	}

	// Run background jobs
	go backgroundJobPastes(time.Second*time.Duration(cfg.DB.CleanupPeriod), db, log)
	go backgroundJobFiles(time.Second*time.Duration(cfg.S3.CleanupPeriod), db, log)

	// Run HTTP server
	log.Info("Run HTTP server on " + cfg.HTTP.Address)
	return handler.Run(log, db, cfg)
}

func main() {
	app := &cli.App{
		Name:    model.SmallName,
		Usage:   "Open source analogue of PasteBin",
		Version: model.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cfg-dir",
				Aliases: []string{"d"},
				Value:   config.DefaultDir,
				Usage:   "directory with Lenpaste config files",
			},
		},
		DefaultCommand:       "help",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run as demon",
				Action: func(ctx *cli.Context) error {
					return run(ctx.String("cfg-dir"))
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fatal(err)
	}
}
