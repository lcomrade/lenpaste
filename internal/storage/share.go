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
	"crypto/rand"
	"math/big"
)

func genTokenCrypto(tokenLen int) (string, error) {
	// Generate token
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	charsLen := int64(len(chars))
	charsLenBig := big.NewInt(charsLen)

	token := ""

	for i := 0; i < tokenLen; i++ {
		randInt, err := rand.Int(rand.Reader, charsLenBig)
		if err != nil {
			return "", err
		}

		token = token + string(chars[randInt.Int64()])
	}

	return token, nil
}
