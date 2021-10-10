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
package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	pasteDir        = "./data/paste"
	pasteFileMod    = 0777
	pasteTextPrefix = ".txt"
	pasteInfoPrefix = ".json"
	pasteNameLength = 8
)

//Type
type PasteInfoType struct {
	CreateTime int64
	DeleteTime int64
	//Title string
	//Syntax string
	//OneUse bool
	//Password string
}

type NewPasteType struct {
	Name string
}

type GetPasteType struct {
	Name string
	Text string
	Info PasteInfoType
}

//Service
func randString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randString := make([]rune, pasteNameLength)
	for i := range randString {
		rand.Seed(time.Now().UnixNano())
		randString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randString)
	//rand.Seed(time.Now().UnixNano())
	//return strconv.Itoa(int(rand.Int63()))
}

func expirParse(expiration string) (int64, error) {
	switch expiration {
	//Minutes
	case "5m":
		return 60 * 5, nil

	case "10m":
		return 60 * 10, nil

	case "20m":
		return 60 * 20, nil

	case "30m":
		return 60 * 30, nil

	case "40m":
		return 60 * 40, nil

	case "50m":
		return 60 * 50, nil

	//Hours
	case "1h":
		return 60 * 60, nil

	case "2h":
		return 60 * 60 * 2, nil

	case "4h":
		return 60 * 60 * 4, nil

	case "12h":
		return 60 * 60 * 12, nil

	//Days
	case "1d":
		return 60 * 60 * 24, nil

	case "2d":
		return 60 * 60 * 24 * 2, nil

	case "3d":
		return 60 * 60 * 24 * 3, nil

	case "4d":
		return 60 * 60 * 24 * 4, nil

	case "5d":
		return 60 * 60 * 24 * 5, nil

	case "6d":
		return 60 * 60 * 24 * 6, nil

	//Weeks
	case "1w":
		return 60 * 60 * 24 * 7, nil

	case "2w":
		return 60 * 60 * 24 * 7 * 2, nil

	case "3w":
		return 60 * 60 * 24 * 7 * 3, nil
	}

	return 0, errors.New("unknown expiration: " + expiration)
}

func genPasteInfo(expir string) ([]byte, error) {
	var infoByte []byte

	//Time
	nowTime := time.Now().Unix()
	expirTime, err := expirParse(expir)
	if err != nil {
		return infoByte, err
	}

	delTime := nowTime + expirTime

	//Struct
	pasteInfo := PasteInfoType{
		CreateTime: nowTime,
		DeleteTime: delTime,
	}

	//Marshal (json)
	infoByte, err = json.Marshal(pasteInfo)
	if err != nil {
		return infoByte, err
	}

	return infoByte, nil
}

func getPasteInfo(name string) (PasteInfoType, error) {
	pasteInfo := PasteInfoType{}

	//File name
	fileName := filepath.Join(pasteDir, name+pasteInfoPrefix)

	//Open file
	file, err := os.Open(fileName)
	if err != nil {
		return pasteInfo, err
	}
	defer file.Close()

	//Decode JSON
	jsonParser := json.NewDecoder(file)

	err = jsonParser.Decode(&pasteInfo)
	if err != nil {
		return pasteInfo, err
	}

	return pasteInfo, nil
}

func isPasteExist(name string) bool {
	//File name
	fileName := filepath.Join(pasteDir, name)
	infoFileName := fileName + pasteInfoPrefix
	textFileName := fileName + pasteTextPrefix

	//Check (INFO)
	_, err := os.Stat(infoFileName)
	if os.IsNotExist(err) == true {
		return false
	}

	//Check (TEXT)
	_, err = os.Stat(textFileName)
	if os.IsNotExist(err) == true {
		return false
	}

	//Return
	return true
}

//Create Paste
func NewPaste(pasteText string, expir string) (NewPasteType, error) {
	var paste NewPasteType

	//Paste name
	pasteName := randString()
	if isPasteExist(pasteName) == true {
		return paste, errors.New("paste with '" + pasteName + "' name exists")
	}

	//File name
	fileName := filepath.Join(pasteDir, pasteName)
	infoFileName := fileName + pasteInfoPrefix
	textFileName := fileName + pasteTextPrefix

	//Paste info
	pasteInfo, err := genPasteInfo(expir)
	if err != nil {
		return paste, err
	}

	//Make paste dir
	err = os.MkdirAll(pasteDir, pasteFileMod)
	if err != nil {
		return paste, err
	}

	//Write (info)
	err = ioutil.WriteFile(infoFileName, pasteInfo, pasteFileMod)
	if err != nil {
		return paste, err
	}

	//Write (text)
	err = ioutil.WriteFile(textFileName, []byte(pasteText), pasteFileMod)
	if err != nil {
		return paste, err
	}

	//Return
	paste = NewPasteType{
		Name: pasteName,
	}

	return paste, nil
}

//Get Paste
func GetPaste(name string) (GetPasteType, error) {
	var pasteInfo PasteInfoType
	var pasteText string
	var paste GetPasteType

	//File name
	fileName := filepath.Join(pasteDir, name)
	textFileName := fileName + pasteTextPrefix

	//Get paste info
	pasteInfo, err := getPasteInfo(name)
	if err != nil {
		return paste, err
	}

	//Get paste text
	pasteTextByte, err := ioutil.ReadFile(textFileName)
	if err != nil {
		return paste, err
	}

	pasteText = string(pasteTextByte)

	//Return
	paste = GetPasteType{
		Name: name,
		Text: pasteText,
		Info: pasteInfo,
	}

	return paste, nil
}

//Delete Paste
func DelPaste(name string) error {
	//File name
	fileName := filepath.Join(pasteDir, name)
	infoFileName := fileName + pasteInfoPrefix
	textFileName := fileName + pasteTextPrefix

	//Remove info
	err := os.Remove(infoFileName)
	if err != nil {
		return err
	}

	//Remove text
	err = os.Remove(textFileName)
	if err != nil {
		return err
	}

	return nil
}

//Deletes expired pastes
func DelExpiredPaste() error {
	//Get file prefix len
	infoPrefixLen := len(pasteInfoPrefix)

	//Get files list
	filesList, err := ioutil.ReadDir(pasteDir)
	if err != nil {
		return err
	}

	//Read files list
	for _, file := range filesList {
		//Skip dirs
		if file.IsDir() == true {
			continue
		}

		//Get file info
		name := file.Name()
		nameLen := len(name)

		//Check file name (info prefix)
		if name[nameLen-infoPrefixLen:] == pasteInfoPrefix {
			//Get paste name
			pasteName := name[0 : nameLen-infoPrefixLen]
			//Get paste info
			paste, err := getPasteInfo(pasteName)
			if err != nil {
				return err
			}

			//Expiry time check
			if paste.DeleteTime <= time.Now().Unix() {
				//Delete expired paste
				err := DelPaste(pasteName)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func init() {
	//Make paste dir (errors are ignored)
	os.MkdirAll(pasteDir, pasteFileMod)
}
