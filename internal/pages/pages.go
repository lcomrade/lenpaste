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
	"../config"
	"../storage"
	"bytes"
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
	Text        string
	ExpiresMin  int64
	ExpiresHour int64
	ExpiresDay  int64
	//Title string
	//Syntax string
	//OneUse bool
	//Password string
}

type errorPageType struct {
	Code  int
	Error string
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
		errorHandler(rw, err, 400)
		return
	}

	//Set Header
	rw.Header().Set("Content-Type", "text/html")

	//Filling the html page template
	tmpl := newDoneTmpl

	err = tmpl.Execute(rw, paste)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

//About API page
func API(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, apiPage)
}

//Server rules page
func Rules(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, rulesPage)
}

//Server version page
func Version(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, versionPage)
}

//Get paste
func GetPaste(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "text/html")

	//Main page
	if req.URL.Path == "/" {
		io.WriteString(rw, mainPage)
		return
	}

	//Get paste name
	name := filepath.Base(req.URL.Path)

	//Get paste
	pasteInfo, err := storage.GetPaste(name)
	if err != nil {
		errorHandler(rw, err, 404)
		return
	}

	//Convert paste info
	deltaTime := pasteInfo.Info.DeleteTime - time.Now().Unix()
	deltaTime = deltaTime / 60
	pastePage := PastePageType{
		Text:        pasteInfo.Text,
		ExpiresMin:  deltaTime % 60,
		ExpiresHour: deltaTime / 60 % 24,
		ExpiresDay:  deltaTime / 60 / 24,
	}

	//Filling the html page template
	tmpl := getTmpl

	err = tmpl.Execute(rw, pastePage)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

//Error page
func errorHandler(rw http.ResponseWriter, err error, code int) {
	//Set Header
	rw.WriteHeader(404)
	rw.Header().Set("Content-Type", "text/html")

	//Get error info
	errorInfo := errorPageType{
		Code:  code,
		Error: err.Error(),
	}

	//Filling the html page template
	tmpl := errorTmpl

	err = tmpl.Execute(rw, errorInfo)
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

func loadRules() (string, error) {
	var rulesHTML string

	//Load HTML template
	tmpl, err := template.ParseFiles(filepath.Join(webDir, "rules.tmpl"))
	if err != nil {
		return rulesHTML, err
	}

	//Read rules file
	rules, err := config.ReadRules()
	if err != nil {
		return rulesHTML, err
	}

	if rules.Exist == false {
		rules.Text = "This server has no rules."
	}

	//Execute template
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, rules)
	if err != nil {
		return rulesHTML, err
	}

	rulesHTML = buf.String()

	return rulesHTML, nil
}

func loadVersion() (string, error) {
	var versionHTML string

	//Load HTML template
	tmpl, err := template.ParseFiles(filepath.Join(webDir, "version.tmpl"))
	if err != nil {
		return versionHTML, err
	}

	//Read version file
	version, err := config.ReadVersion()
	if err != nil {
		return versionHTML, err
	}

	//Execute template
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, version)
	if err != nil {
		return versionHTML, err
	}

	versionHTML = buf.String()

	return versionHTML, nil
}

var styleCSS string
var mainPage string
var apiPage string
var rulesPage string
var versionPage string
var newPage string
var newDoneTmpl *template.Template
var getTmpl *template.Template
var errorTmpl *template.Template

func Load() error {
	//Style
	styleCSSByte, err := loadFile(filepath.Join(webDir, "style.css"))
	if err != nil {
		return err
	}

	styleCSS = string(styleCSSByte)

	//Main page
	mainPageByte, err := loadFile(filepath.Join(webDir, "main.html"))
	if err != nil {
		return err
	}

	mainPage = string(mainPageByte)

	//About API page
	apiPageByte, err := loadFile(filepath.Join(webDir, "api.html"))
	if err != nil {
		return err
	}

	apiPage = string(apiPageByte)

	//Rules page
	rulesPageLoad, err := loadRules()
	if err != nil {
		return err
	}

	rulesPage = rulesPageLoad

	//Version page
	versionPageLoad, err := loadVersion()
	if err != nil {
		return err
	}

	versionPage = versionPageLoad

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

	//Error tmpl
	errorTmplLoad, err := template.ParseFiles(filepath.Join(webDir, "error.tmpl"))
	if err != nil {
		return err
	}

	errorTmpl = errorTmplLoad

	return nil
}
