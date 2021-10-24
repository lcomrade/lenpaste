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
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

const (
	configPath  = "./data/config.json"
	aboutPath   = "./data/about.txt"
	rulesPath   = "./data/rules.txt"
	versionPath = "./version.json"
)

//Config file TYPE
type Config struct {
	HTTP    ConfigHTTP
	Storage ConfigStorage
	Logs    ConfigLogs
}

type ConfigHTTP struct {
	Listen  string
	UseTLS  bool
	SSLCert string
	SSLKey  string
}

type ConfigStorage struct {
	CleanJobPeriod time.Duration
}

type ConfigLogs struct {
	SaveErr    bool
	SaveInfo   bool
	SaveJob    bool
	RotateLogs bool
	MaxLogSize int64
}

//Config DEFAULT
var defaultCfg = Config{
	HTTP: ConfigHTTP{
		Listen:  ":8000",
		UseTLS:  false,
		SSLCert: "./data/fullchain.pem",
		SSLKey:  "./data/privkey.pem",
	},
	Storage: ConfigStorage{
		CleanJobPeriod: 10 * time.Minute, //10 minutes
	},
	Logs: ConfigLogs{
		SaveErr:    true,
		SaveInfo:   true,
		SaveJob:    true,
		RotateLogs: true,
		MaxLogSize: 1000000, //1 MB
	},
}

// Get config file
func ReadConfig() (Config, error) {
	//Set default values
	config := defaultCfg

	//Read config
	file, err := os.Open(configPath)
	if err != nil {
		//If the config is missing
		if os.IsNotExist(err) == true {
			return config, nil

			//If another error
		} else {
			return config, err
		}
	}

	//Decode config file
	parser := json.NewDecoder(file)
	err = parser.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

//Get text file
type TextType struct {
	Exist bool
	Text  string
}

var TextDefault = TextType{
	Exist: false,
	Text:  "",
}

func readText(path string) (TextType, error) {
	out := TextDefault

	//Open text file
	file, err := os.Open(path)
	//If text file is missing
	if err != nil {
		if os.IsNotExist(err) == true {
			return out, nil

			//If another error
		} else {
			return out, err
		}
	}
	defer file.Close()

	//Read text file
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return out, err
	}

	//Byte to string
	out.Text = string(fileByte)
	out.Exist = true

	//Return
	return out, nil
}

//Get about text
func ReadAbout() (TextType, error) {
	out, err := readText(aboutPath)
	return out, err
}

//Get rules text
func ReadRules() (TextType, error) {
	out, err := readText(rulesPath)
	return out, err
}

//Get version
type VersionType struct {
	Version   string
	GitTag    string
	GitCommit string
	BuildDate string
}

var VersionDefault = VersionType{
	Version:   "unknown",
	GitTag:    "unknown",
	GitCommit: "unknown",
	BuildDate: "unknown",
}

func ReadVersion() (VersionType, error) {
	//Set default values
	version := VersionDefault

	//Read version file
	file, err := os.Open(versionPath)
	if err != nil {
		//If the version is missing
		if os.IsNotExist(err) == true {
			return version, nil

			//If another error
		} else {
			return version, err
		}
	}

	//Decode version file
	parser := json.NewDecoder(file)
	err = parser.Decode(&version)
	if err != nil {
		return version, err
	}

	//Return
	return version, nil
}
