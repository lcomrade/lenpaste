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

package netshare

import (
	"errors"
)

var (
	ErrBadRequest       = errors.New("Bad Request")           // 400
	ErrNotFound         = errors.New("Not Found")             // 404
	ErrMethodNotAllowed = errors.New("Method Not Allowed")    // 405
	ErrInternal         = errors.New("Internal Server Error") // 500
)
