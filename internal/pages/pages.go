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
	Title       string
	Text        string
	OneUse      bool
	ExpiresMin  int64
	ExpiresHour int64
	ExpiresDay  int64
	//Syntax string
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

	io.WriteString(rw, pages.StyleCSS)
}

//New paste
func NewPaste(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, pages.New)
}

//New paste done
type NewPasteType struct {
	Name string
	Host string
}

func NewPasteDone(rw http.ResponseWriter, req *http.Request) {
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

	//Total paste info
	pasteTotal := NewPasteType{
		Name: paste.Name,
		Host: req.Host,
	}

	//Set Header
	rw.Header().Set("Content-Type", "text/html")

	//Filling the html page template
	tmpl := pages.NewDoneTmpl

	err = tmpl.Execute(rw, pasteTotal)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

//About API page
func API(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, pages.API)
}

//Server rules page
func Rules(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, pages.Rules)
}

//Server version page
func Version(rw http.ResponseWriter, req *http.Request) {
	//Return response
	rw.Header().Set("Content-Type", "text/html")

	io.WriteString(rw, pages.Version)
}

//Get paste
func GetPaste(rw http.ResponseWriter, req *http.Request) {
	//Set Header
	rw.Header().Set("Content-Type", "text/html")

	//Main page
	if req.URL.Path == "/" {
		io.WriteString(rw, pages.Main)
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
		Title:       pasteInfo.Info.Title,
		Text:        pasteInfo.Text,
		OneUse:      pasteInfo.Info.OneUse,
		ExpiresMin:  deltaTime % 60,
		ExpiresHour: deltaTime / 60 % 24,
		ExpiresDay:  deltaTime / 60 / 24,
	}

	//Filling the html page template
	tmpl := pages.GetTmpl

	err = tmpl.Execute(rw, pastePage)
	if err != nil {
		errorHandler(rw, err, 400)
		return
	}
}

//Error page
func errorHandler(rw http.ResponseWriter, err error, code int) {
	//Set Header
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "text/html")

	//Get error info
	errorInfo := errorPageType{
		Code:  code,
		Error: err.Error(),
	}

	//Filling the html page template
	tmpl := pages.ErrorTmpl

	err = tmpl.Execute(rw, errorInfo)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

//Load pages (init)
func loadFile(path string) (string, error) {
	var out string

	//Open file
	file, err := os.Open(path)
	if err != nil {
		return out, err
	}
	defer file.Close()

	//Read file
	fileByte, err := ioutil.ReadAll(file)
	out = string(fileByte)
	if err != nil {
		return out, err
	}

	return out, nil
}

func loadMain() (string, error) {
	var mainHTML string

	//Load HTML template
	tmpl, err := template.ParseFiles(filepath.Join(webDir, "main.tmpl"))
	if err != nil {
		return mainHTML, err
	}

	//Read about file
	about, err := config.ReadAbout()
	if err != nil {
		return mainHTML, err
	}

	//Execute template
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, about)
	if err != nil {
		return mainHTML, err
	}

	mainHTML = buf.String()

	return mainHTML, nil
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

type pagesType struct {
	StyleCSS    string
	Main        string
	API         string
	Rules       string
	Version     string
	New         string
	NewDoneTmpl *template.Template
	GetTmpl     *template.Template
	ErrorTmpl   *template.Template
}

var pages pagesType

func Load() error {
	var err error

	//Style
	pages.StyleCSS, err = loadFile(filepath.Join(webDir, "style.css"))
	if err != nil {
		return err
	}

	//Main page
	pages.Main, err = loadMain()
	if err != nil {
		return err
	}

	//About API page
	pages.API, err = loadFile(filepath.Join(webDir, "api.html"))
	if err != nil {
		return err
	}

	//Rules page
	pages.Rules, err = loadRules()
	if err != nil {
		return err
	}

	//Version page
	pages.Version, err = loadVersion()
	if err != nil {
		return err
	}

	//New page
	pages.New, err = loadFile(filepath.Join(webDir, "new.html"))
	if err != nil {
		return err
	}

	//New done tmpl
	pages.NewDoneTmpl, err = template.ParseFiles(filepath.Join(webDir, "new_done.tmpl"))
	if err != nil {
		return err
	}

	//Get tmpl
	pages.GetTmpl, err = template.ParseFiles(filepath.Join(webDir, "get.tmpl"))
	if err != nil {
		return err
	}

	//Error tmpl
	pages.ErrorTmpl, err = template.ParseFiles(filepath.Join(webDir, "error.tmpl"))
	if err != nil {
		return err
	}

	return nil
}
