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

package lenpasswd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Data map[string]string

func LoadFile(path string) (Data, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("lenpasswd: " + err.Error())
	}
	defer file.Close()

	// Read file
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("lenpasswd: " + err.Error())
	}

	// Convert []byte to string
	fileTxt := bytes.NewBuffer(fileByte).String()

	// Parse file
	data := make(Data)
	for i, line := range strings.Split(fileTxt, "\n") {
		if line == "" {
			continue
		}

		lineSplit := strings.Split(line, ":")
		if len(lineSplit) != 2 {
			return nil, errors.New("lenpasswd: error in line " + strconv.Itoa(i))
		}

		user := lineSplit[0]
		pass := lineSplit[1]

		_, exist := data[user]
		if exist == true {
			return nil, errors.New("lenpasswd: overriding user " + user + " in line " + strconv.Itoa(i))
		}

		data[user] = pass
	}

	return data, nil
}

func (data Data) Check(user string, pass string) bool {
	truePass, exist := data[user]
	if exist == false {
		return false
	}

	if pass != truePass {
		return false
	}

	return true
}

func LoadAndCheck(path string, user string, pass string) (bool, error) {
	data, err := LoadFile(path)
	if err != nil {
		return false, err
	}

	return data.Check(user, pass), nil
}
