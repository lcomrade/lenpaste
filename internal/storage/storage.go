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
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFoundID = errors.New("db: could not find ID")
)

type DB struct {
	pool *sql.DB
}

func NewPool(driverName string, dataSourceName string, maxOpenConns int, maxIdleConns int) (DB, error) {
	var db DB
	var err error

	db.pool, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return db, err
	}

	db.pool.SetMaxOpenConns(maxOpenConns)
	db.pool.SetMaxIdleConns(maxIdleConns)

	return db, nil
}

func (db DB) Close() error {
	return db.pool.Close()
}

func InitDB(driverName string, dataSourceName string) error {
	// Open DB
	db, err := NewPool(driverName, dataSourceName, 1, 0)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables
	_, err = db.pool.Exec(`
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
		return err
	}

	// Crutch for SQLite3
	if driverName == "sqlite3" {
		_, err = db.pool.Exec(`ALTER TABLE pastes ADD COLUMN author       TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			if err.Error() != "duplicate column name: author" {
				return err
			}
		}

		_, err = db.pool.Exec(`ALTER TABLE pastes ADD COLUMN author_email TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			if err.Error() != "duplicate column name: author_email" {
				return err
			}
		}

		_, err = db.pool.Exec(`ALTER TABLE pastes ADD COLUMN author_url TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			if err.Error() != "duplicate column name: author_url" {
				return err
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
			return err
		}
	}

	return nil
}
