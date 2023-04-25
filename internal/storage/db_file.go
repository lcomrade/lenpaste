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
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (db *DB) FileUpload(file model.File) (string, int64, int64, string, error) {
	var err error

	// Generate ID
	file.ID, err = genTokenCrypto(8)
	if err != nil {
		return "", 0, 0, "", errors.New("storage: upload file: " + err.Error())
	}

	fileS3Key := file.S3Key()

	file.URL, err = url.JoinPath(db.cfg.S3.URL, fileS3Key)
	if err != nil {
		return "", 0, 0, "", errors.New("storage: upload file: " + err.Error())
	}

	// Set file create time
	file.CreateTime = time.Now().Unix()

	// Check delete time
	if file.DeleteTime < 0 {
		file.DeleteTime = 0
	}

	// Generate S3 pre-signed URL to upload file
	preSignedURL, err := s3.NewPresignClient(db.s3).PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(db.cfg.S3.Bucket),
			Key:    aws.String(fileS3Key),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = 60 * time.Minute
		},
	)
	if err != nil {
		return "", 0, 0, "", errors.New("storage: upload file: " + err.Error())
	}

	// Write to DB
	_, err = db.pool.Exec(
		`INSERT INTO files (
			id, title,
			create_time, delete_time,
			filename, url,
			author, author_email, author_url
		) VALUES (
			$1, $2,
			$3, $4,
			$5, $6,
			$7, $8, $9
		);`,
		file.ID, file.Title,
		file.CreateTime, file.DeleteTime,
		file.Filename, file.URL,
		file.Author, file.AuthorEmail, file.AuthorURL,
	)
	if err != nil {
		return "", 0, 0, "", errors.New("storage: upload file: " + err.Error())
	}

	return file.ID, file.CreateTime, file.DeleteTime, preSignedURL.URL, nil
}

func (db *DB) FileGet(id string) (model.File, error) {
	var file model.File

	// Make query
	row := db.pool.QueryRow(
		`SELECT
			id, title,
			create_time, delete_time,
			filename, url,
			author, author_email, author_url
		FROM files WHERE id = $1;`,
		id,
	)

	// Read query
	err := row.Scan(
		&file.ID, &file.Title,
		&file.CreateTime, &file.DeleteTime,
		&file.Filename, &file.URL,
		&file.Author, &file.AuthorEmail, &file.AuthorURL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.File{}, ErrNotFoundID
		}

		return model.File{}, errors.New("storage: get file: " + err.Error())
	}

	// Check file expiration
	if file.DeleteTime < time.Now().Unix() && file.DeleteTime > 0 {
		// Delete expired file
		_, err = db.pool.Exec(
			`DELETE FROM files WHERE id = $1;`,
			file.ID,
		)
		if err != nil {
			return model.File{}, errors.New("storage: get file: " + err.Error())
		}

		// Return ErrNotFound
		return model.File{}, ErrNotFoundID
	}

	return file, nil
}

func (db *DB) FileCleanup() (int64, int64, error) {
	var expired, notFinished int64

	// Remove expired files list
	{
		timeNowUnix := time.Now().Unix()

		rows, err := db.pool.Query(
			`SELECT id, filename FROM files WHERE (delete_time < $1) AND (delete_time > 0) AND (upload_finished = true);`,
			timeNowUnix,
		)
		if err != nil {
			return 0, 0, errors.New("storage: cleanup uploaded files (expiration): " + err.Error())
		}

		// Prepare S3 objects list
		var s3Objects []s3Types.ObjectIdentifier
		for rows.Next() {
			var id, filename string
			err := rows.Scan(&id, &filename)
			if err != nil {
				return 0, 0, errors.New("storage: cleanup uploaded files (expiration): " + err.Error())
			}

			s3Objects = append(
				s3Objects,
				s3Types.ObjectIdentifier{
					Key: aws.String(model.FileS3Key(id, filename)),
				},
			)
		}

		// Remove from S3 storage
		_, err = db.s3.DeleteObjects(
			context.TODO(),
			&s3.DeleteObjectsInput{
				Bucket: aws.String(db.cfg.S3.Bucket),
				Delete: &s3Types.Delete{Objects: s3Objects},
			},
		)
		if err != nil {
			return 0, 0, errors.New("storage: cleanup uploaded files (expiration): " + err.Error())
		}

		// Remove from DB
		_, err = db.pool.Exec(
			`DELETE FROM files WHERE (delete_time < $1) AND (delete_time > 0) AND (upload_finished = true);`,
			timeNowUnix,
		)
		if err != nil {
			return 0, 0, errors.New("storage: cleanup uploaded files (expiration): " + err.Error())
		}

		expired = int64(len(s3Objects))
	}

	// Remove not finished uploads
	{
		result, err := db.pool.Exec(
			`DELETE FROM files WHERE (upload_finished = false) AND (create_time > $1);`,
			time.Now().Unix(),
		)
		if err != nil {
			return 0, 0, errors.New("storage: cleanup uploaded files (not finished): " + err.Error())
		}

		notFinished, err = result.RowsAffected()
		if err != nil {
			return 0, 0, errors.New("storage: cleanup uploaded files (not finished): " + err.Error())
		}
	}

	return expired, notFinished, nil
}
