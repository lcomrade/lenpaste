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

package storage

import (
	"database/sql"
	"time"
)

type Paste struct {
	ID         string `json:"id"` // Ignored when creating
	Title      string `json:"title"`
	Body       string `json:"body"`
	CreateTime int64  `json:"createTime"` // Ignored when creating
	DeleteTime int64  `json:"deleteTime"`
	OneUse     bool   `json:"oneUse"`
	Syntax     string `json:"syntax"`
	//Password string `json:"password"`
}

func (dbInfo DB) PasteAdd(paste Paste) (Paste, error) {
	// Open DB
	db, err := dbInfo.openDB()
	if err != nil {
		return paste, err
	}
	defer db.Close()

	// Generate ID
	paste.ID, err = genTokenCrypto(8)
	if err != nil {
		return paste, err
	}

	// Set paste create time
	paste.CreateTime = time.Now().Unix()

	// Check delete time
	if paste.DeleteTime < 0 {
		paste.DeleteTime = 0
	}

	// Add
	_, err = db.Exec(
		`INSERT INTO "pastes" ("id", "title", "body", "syntax", "create_time", "delete_time", "one_use") VALUES (?, ?, ?, ?, ?, ?, ?)`,
		paste.ID, paste.Title, paste.Body, paste.Syntax, paste.CreateTime, paste.DeleteTime, paste.OneUse,
	)
	if err != nil {
		return paste, err
	}

	return paste, nil
}

func (dbInfo DB) PasteDelete(id string) error {
	// Open DB
	db, err := dbInfo.openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete
	result, err := db.Exec(
		`DELETE FROM "pastes" WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	// Check result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFoundID
	}

	return nil
}

func (dbInfo DB) PasteGet(id string) (Paste, error) {
	var paste Paste

	// Open DB
	db, err := dbInfo.openDB()
	if err != nil {
		return paste, err
	}
	defer db.Close()

	// Make query
	row := db.QueryRow(
		`SELECT "id", "title", "body", "syntax", "create_time", "delete_time", "one_use" FROM "pastes" WHERE "id" = ?`,
		id,
	)

	// Read query
	err = row.Scan(&paste.ID, &paste.Title, &paste.Body, &paste.Syntax, &paste.CreateTime, &paste.DeleteTime, &paste.OneUse)
	if err != nil {
		if err == sql.ErrNoRows {
			return paste, ErrNotFoundID
		}

		return paste, err
	}

	// Check paste expiration
	if paste.DeleteTime < time.Now().Unix() && paste.DeleteTime > 0 {
		// Delete expired paste
		_, err = db.Exec(
			`DELETE FROM "pastes" WHERE id = ?`,
			paste.ID,
		)
		if err != nil {
			return Paste{}, err
		}

		// Return ErrNotFound
		return Paste{}, ErrNotFoundID
	}

	return paste, nil
}

func (dbInfo DB) PasteGetList() ([]Paste, error) {
	var pastes []Paste

	// Open DB
	db, err := dbInfo.openDB()
	if err != nil {
		return pastes, err
	}
	defer db.Close()

	// Delete expired paste
	_, err = db.Exec(
		`DELETE FROM "pastes" WHERE delete_time < ? AND delete_time > 0`,
		time.Now().Unix(),
	)
	if err != nil {
		return pastes, err
	}

	// Make query to get paste list
	rows, err := db.Query(
		`SELECT "id", "title", "body", "syntax", "create_time", "delete_time", "one_use" FROM "pastes"`,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return pastes, ErrNotFoundID
		}

		return pastes, err
	}

	// Read query
	for rows.Next() {
		var paste Paste

		err = rows.Scan(&paste.ID, &paste.Title, &paste.Body, &paste.Syntax, &paste.CreateTime, &paste.DeleteTime, &paste.OneUse)
		if err != nil {
			return pastes, err
		}

		pastes = append(pastes, paste)
	}

	return pastes, nil
}

func (dbInfo DB) PasteDeleteExpired() (int64, error) {
	// Open DB
	db, err := dbInfo.openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Delete
	result, err := db.Exec(
		`DELETE FROM "pastes" WHERE delete_time < ? AND delete_time > 0`,
		time.Now().Unix(),
	)
	if err != nil {
		return 0, err
	}

	// Check result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}
