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
package pages

import (
	"../storage"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	webDir = "./web"
)

type PastePageType struct {
	Text       string
	ExpiresMin int64
	ExpiresDay int64
	//Title string
	//Syntax string
	//OneUse bool
	//Password string
}

//Style
func Style(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/css")

	io.WriteString(rw, styleCSS)
}

//New paste
func NewPaste(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, newPage)
}

//New paste done
func NewPasteDone(rw http.ResponseWriter, req *http.Request) {
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

	//Set Header
	rw.Header().Set("Content-Type", "text/html")

	//Filling the html page template
	tmpl := newDoneTmpl

	err = tmpl.Execute(rw, paste)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

//Get paste
func GetPaste(rw http.ResponseWriter, req *http.Request) {
	//Get paste name
	name := filepath.Base(req.URL.Path)

	//Get paste
	pasteInfo, err := storage.GetPaste(name)
	if err != nil {
		http.Error(rw, err.Error(), 404)
		return
	}

	//Convert paste info
	deltaTime := pasteInfo.Info.DeleteTime - time.Now().Unix()
	deltaTime = deltaTime / 60
	pastePage := PastePageType{
		Text:       pasteInfo.Text,
		ExpiresMin: deltaTime % 60,
		ExpiresDay: deltaTime / 60 / 24,
	}

	//Set Header
	rw.Header().Set("Content-Type", "text/html")

	//Filling the html page template
	tmpl := getTmpl

	err = tmpl.Execute(rw, pastePage)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

//Load pages (init)
func loadFile(path string) ([]byte, error) {
	var fileByte []byte

	//Open file
	file, err := os.Open(path)
	if err != nil {
		return fileByte, err
	}
	defer file.Close()

	//Read file
	fileByte, err = ioutil.ReadAll(file)
	if err != nil {
		return fileByte, err
	}

	return fileByte, nil
}

var styleCSS string
var newPage string
var newDoneTmpl *template.Template
var getTmpl *template.Template

func Load() error {
	//Style
	styleCSSByte, err := loadFile(filepath.Join(webDir, "style.css"))
	if err != nil {
		return err
	}

	styleCSS = string(styleCSSByte)

	//New page
	newPageByte, err := loadFile(filepath.Join(webDir, "new.html"))
	if err != nil {
		return err
	}

	newPage = string(newPageByte)

	//New done tmpl
	newDoneTmplLoad, err := template.ParseFiles(filepath.Join(webDir, "new_done.tmpl"))
	if err != nil {
		return err
	}

	newDoneTmpl = newDoneTmplLoad

	//Get tmpl
	getTmplLoad, err := template.ParseFiles(filepath.Join(webDir, "get.tmpl"))
	if err != nil {
		return err
	}

	getTmpl = getTmplLoad

	return nil
}
