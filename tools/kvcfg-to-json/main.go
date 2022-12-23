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
