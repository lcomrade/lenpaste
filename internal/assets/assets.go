/*
   Copyright 2021 Leonid Maslakov

   License: GPL-3.0-or-later

   This file is part of Lenpaste.

   Lenpaste is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Lenpaste is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Lenpaste.  If not, see <https://www.gnu.org/licenses/>.
*/
package assets

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	robotsTxtPath = "./data/robots.txt"
)

//robots.txt
func RobotsTxt(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/plain")

	io.WriteString(rw, robotsTxt)
}

//Load robots.txt
func loadRobotsTxt() string {
	robotsTxt := `
User-agent: *
Disallow: /`

	//Open file
	file, err := os.Open(robotsTxtPath)
	if err != nil {
		return robotsTxt
	}
	defer file.Close()

	//Read file
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return robotsTxt
	}

	return string(fileByte)
}

var robotsTxt string

//Load
func Load() error {
	robotsTxt = loadRobotsTxt()

	return nil
}
