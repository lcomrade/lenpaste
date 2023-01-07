// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stdout, "Usage:", os.Args[0], "[SRC] [DST]")
		os.Exit(1)
	}

	src := os.Args[1]
	dst := os.Args[2]

	// Read key-value config
	fileByte, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}

	cfg, err := readKVCfg(string(fileByte))
	if err != nil {
		panic(err)
	}

	// Save config as JSON
	cfgRaw, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		panic(err)
	}

	cfgRaw = append(cfgRaw, byte('\n'))

	err = os.WriteFile(dst, cfgRaw, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
