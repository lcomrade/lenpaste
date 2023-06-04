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

import "errors"

type Error struct {
	Code    int               `json:"code"`
	Text    string            `json:"error"`
	RealErr string            `json:"-"`
	Header  map[string]string `json:"-"`
}

func (e *Error) Error() string {
	return e.RealErr
}

func NewError(code int, text string) error {
	return &Error{
		Code:    code,
		Text:    text,
		RealErr: text,
	}
}
func NewErrorWithHeader(code int, text string, header map[string]string) error {
	return &Error{
		Code:    code,
		Text:    text,
		RealErr: text,
		Header:  header,
	}
}

func ParseError(e error) Error {
	var resp *Error

	if !errors.As(e, &resp) {
		resp = &Error{
			Code:    500,
			Text:    "Internal Server Error",
			RealErr: e.Error(),
		}
	}

	return *resp
}
