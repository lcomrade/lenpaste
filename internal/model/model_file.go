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

package model

type File struct {
	ID         string `json:"id"` // Ignored when creating
	Title      string `json:"title"`
	CreateTime int64  `json:"createTime"` // Ignored when creating
	DeleteTime int64  `json:"deleteTime"`

	Filename string `json:"filename"`
	URL      string `json:"url"`

	Author      string `json:"author"`
	AuthorEmail string `json:"authorEmail"`
	AuthorURL   string `json:"authorURL"`
}

func (f *File) S3Key() string {
	return FileS3Key(f.ID, f.Filename)
}

func FileS3Key(id, filename string) string {
	return id + "_" + filename
}
