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
package api

import (
	"../config"
	"../storage"
	"encoding/json"
	"net/http"
	"path/filepath"
)

type errorType struct {
	Code  int
	Error string
}

func errorHandler(rw http.ResponseWriter, err error, code int) {
	//Set Header
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json")

	//Get error info
	errorInfo := errorType{
		Code:  code,
		Error: err.Error(),
	}

	//Return response
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&errorInfo)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

func NewPaste(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "application/json")

	//Get form data
	req.ParseForm()

	text := req.Form.Get("text")
	expiration := req.Form.Get("expiration")
	title := req.Form.Get("title")

	oneUse := false
	if req.Form.Get("oneUse") == "true" {
		oneUse = true
	}

	//Create paste
	paste, err := storage.NewPaste(text, expiration, oneUse, title)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}

	//Return response
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&paste)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

func GetPaste(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "application/json")

	//Get paste name
	name := filepath.Base(req.URL.Path)

	//Get paste
	paste, err := storage.GetPaste(name)
	if err != nil {
		errorHandler(rw, err, 404)
		return
	}

	//Return response
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&paste)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

func GetAbout(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "application/json")

	//Get about
	about, err := config.ReadAbout()
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}

	//Return response
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&about)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

func GetRules(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "application/json")

	//Get rules
	rules, err := config.ReadRules()
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}

	//Return response
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&rules)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

func GetVersion(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "application/json")

	//Get version
	version, err := config.ReadVersion()
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}

	//Return response
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&version)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}
