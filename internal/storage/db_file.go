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
	"errors"
	"net/url"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (db *DB) UploadFile(file model.File) (string, int64, int64, string, error) {
	var err error

	// Generate ID
	file.ID, err = genTokenCrypto(8)
	if err != nil {
		return "", 0, 0, "", errors.New("storage: upload file: " + err.Error())
	}

	fileS3Key := file.ID + "_" + file.Filename

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
		`INSERT INTO pastes (
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
