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
	"os"
)

const (
	configPath = "./data/config.json"
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
func Read() (Config, error) {
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
