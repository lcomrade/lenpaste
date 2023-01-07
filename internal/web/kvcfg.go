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

package web

import (
	"errors"
	"strconv"
	"strings"
)

func readKVCfg(data string) (map[string]string, error) {
	out := make(map[string]string)

	dataSplit := strings.Split(data, "\n")
	dataSplitLen := len(dataSplit)

	for num := 0; num < dataSplitLen; num++ {
		str := strings.TrimSpace(dataSplit[num])

		if str == "" || strings.HasPrefix(str, "//") {
			continue
		}

		strSplit := strings.SplitN(str, "=", 2)
		if len(strSplit) != 2 {
			return out, errors.New("error in line " + strconv.Itoa(num+1) + ": expected '=' delimiter")
		}

		key := strings.TrimSpace(strSplit[0])
		val := strings.TrimSpace(strSplit[1])
		val, isMultiline := multilineCheck(val)

		if isMultiline {
			num = num + 1
			for ; num < dataSplitLen; num++ {
				strPlus := strings.TrimSpace(dataSplit[num])
				strPlus, isMultilinePlus := multilineCheck(strPlus)
				val = val + strPlus

				if isMultilinePlus == false {
					break
				}
			}
		}

		_, exist := out[key]
		if exist {
			return out, errors.New("duplicate key: " + key)
		}

		out[key] = val
	}

	return out, nil
}

func multilineCheck(s string) (string, bool) {
	sLen := len(s)

	if sLen > 0 && s[sLen-1] == '\\' {
		if sLen > 1 && s[sLen-2] == '\\' {
			return s[:sLen-1], false
		}

		return s[:sLen-1], true
	}

	return s, false
}
