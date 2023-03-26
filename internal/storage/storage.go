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

package storage

import (
	"database/sql"
	"errors"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFoundID = errors.New("db: could not find ID")
)

type DB struct {
	cfg  *config.Config
	pool *sql.DB
}

func Open(cfg *config.Config) (*DB, error) {
	// Open SQL DB
	var err error
	db := DB{
		cfg: cfg,
	}

	db.pool, err = sql.Open(db.cfg.DB.Driver, db.cfg.DB.Source)
	if err != nil {
		return nil, errors.New("storage: open: " + err.Error())
	}

	// Setup SQL DB
	db.pool.SetMaxOpenConns(db.cfg.DB.MaxOpenConns)
	db.pool.SetMaxIdleConns(db.cfg.DB.MaxIdleConns)
	db.pool.SetConnMaxLifetime(time.Duration(cfg.DB.ConnMaxLifetime) * time.Second)

	return &db, nil
}

func (db DB) Close() error {
	return db.pool.Close()
}

func (db *DB) InitDB() error {
	// Create tables
	_, err := db.pool.Exec(`
		CREATE TABLE IF NOT EXISTS pastes (
			id          TEXT    PRIMARY KEY,
			title       TEXT    NOT NULL,
			body        TEXT    NOT NULL,
			syntax      TEXT    NOT NULL,
			create_time INTEGER NOT NULL,
			delete_time INTEGER NOT NULL,
			one_use     BOOL    NOT NULL
		);
	`)
	if err != nil {
		return errors.New("storage: init: " + err.Error())
	}

	// Crutch for SQLite3
	if db.cfg.DB.Driver == "sqlite3" {
		_, err = db.pool.Exec(`ALTER TABLE pastes ADD COLUMN author       TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			if err.Error() != "duplicate column name: author" {
				return errors.New("storage: init: " + err.Error())
			}
		}

		_, err = db.pool.Exec(`ALTER TABLE pastes ADD COLUMN author_email TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			if err.Error() != "duplicate column name: author_email" {
				return errors.New("storage: init: " + err.Error())
			}
		}

		_, err = db.pool.Exec(`ALTER TABLE pastes ADD COLUMN author_url TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			if err.Error() != "duplicate column name: author_url" {
				return errors.New("storage: init: " + err.Error())
			}
		}

		// Normal SQL for all other DBs
	} else {
		_, err = db.pool.Exec(`
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS author       TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS author_email TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS author_url   TEXT NOT NULL DEFAULT '';
		`)
		if err != nil {
			return errors.New("storage: init: " + err.Error())
		}
	}

	return nil
}
