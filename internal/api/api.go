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

func NewPaste(rw http.ResponseWriter, req *http.Request) {
	//Get form data
	req.ParseForm()
	text := req.Form.Get("text")
	expiration := req.Form.Get("expiration")

	//Create paste
	paste, err := storage.NewPaste(text, expiration)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	//Return response
	rw.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&paste)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func GetPaste(rw http.ResponseWriter, req *http.Request) {
	//Get paste name
	name := filepath.Base(req.URL.Path)

	//Get paste
	paste, err := storage.GetPaste(name)
	if err != nil {
		http.Error(rw, err.Error(), 404)
		return
	}

	//Return response
	rw.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&paste)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func GetAbout(rw http.ResponseWriter, req *http.Request) {
	//Get about
	about, err := config.ReadAbout()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	//Return response
	rw.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&about)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func GetRules(rw http.ResponseWriter, req *http.Request) {
	//Get rules
	rules, err := config.ReadRules()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	//Return response
	rw.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&rules)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func GetVersion(rw http.ResponseWriter, req *http.Request) {
	//Get version
	version, err := config.ReadVersion()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	//Return response
	rw.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&version)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}
