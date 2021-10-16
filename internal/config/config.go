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
)

const (
	configPath  = "./data/config.json"
	aboutPath   = "./data/about.txt"
	rulesPath   = "./data/rules.txt"
	versionPath = "./version.json"
)

//Config file TYPE
type Config struct {
	HTTP ConfigHTTP
}

type ConfigHTTP struct {
	Listen  string
	UseTLS  bool
	SSLCert string
	SSLKey  string
}

//Config DEFAULT
var defaultCfg = Config{
	HTTP: ConfigHTTP{
		Listen:  ":8000",
		UseTLS:  false,
		SSLCert: "./data/fullchain.pem",
		SSLKey:  "./data/privkey.pem",
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

//Get about text
type AboutType struct {
	Exist bool
	Text  string
}

var AboutDefault = AboutType{
	Exist: false,
	Text:  "",
}

func ReadAbout() (AboutType, error) {
	about := AboutDefault

	//Open about file
	file, err := os.Open(aboutPath)
	//If the about file is missing
	if err != nil {
		if os.IsNotExist(err) == true {
			return about, nil

			//If another error
		} else {
			return about, err
		}
	}
	defer file.Close()

	//Read rules file
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return about, err
	}

	//Byte to string
	about.Text = string(fileByte)
	about.Exist = true

	//Return
	return about, nil
}

//Get rules text
type RulesType struct {
	Exist bool
	Text  string
}

var RulesDefault = RulesType{
	Exist: false,
	Text:  "",
}

func ReadRules() (RulesType, error) {
	rules := RulesDefault

	//Open rules file
	file, err := os.Open(rulesPath)
	//If the rules file is missing
	if err != nil {
		if os.IsNotExist(err) == true {
			return rules, nil

			//If another error
		} else {
			return rules, err
		}
	}
	defer file.Close()

	//Read rules file
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return rules, err
	}

	//Byte to string
	rules.Text = string(fileByte)
	rules.Exist = true

	//Return
	return rules, nil
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
