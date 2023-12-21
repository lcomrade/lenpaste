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

package cli

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func exitOnError(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}

type variable struct {
	name        string
	cliFlagName string

	preHook func(string) (string, error)

	value        interface{}
	valueDefault string
	required     bool
	usage        string
}

type CLI struct {
	version string

	vars []variable
}

type FlagOptions struct {
	Required bool
	PreHook  func(string) (string, error)
}

func New(version string) *CLI {
	return &CLI{
		version: version,

		vars: []variable{},
	}
}

func (c *CLI) addVar(name string, value interface{}, defValue string, usage string, opts *FlagOptions) {
	if name == "" {
		panic("cli: add variable: variable name could not be empty")
	}

	if usage == "" {
		panic("cli: flag \"" + name + "\" has empty \"usage\" field")
	}

	if opts == nil {
		opts = &FlagOptions{}
	}

	c.vars = append(c.vars, variable{
		name:        name,
		cliFlagName: "-" + name,

		preHook: opts.PreHook,

		value:        value,
		valueDefault: defValue,
		required:     opts.Required,
		usage:        usage,
	})
}

func (c *CLI) AddStringVar(name, defValue string, usage string, opts *FlagOptions) *string {
	if opts != nil {
		if opts.PreHook != nil {
			var err error
			defValue, err = opts.PreHook(defValue)
			if err != nil {
				panic("cli: add duration variable \"" + name + "\": " + err.Error())
			}
		}
	}

	val := &defValue
	c.addVar(name, val, defValue, usage, opts)
	return val
}

func (c *CLI) AddBoolVar(name string, usage string) *bool {
	valVar := false
	val := &valVar
	c.addVar(name, val, "", usage, nil)
	return val
}

func (c *CLI) AddIntVar(name string, defValue int, usage string, opts *FlagOptions) *int {
	val := &defValue
	c.addVar(name, val, strconv.Itoa(defValue), usage, opts)
	return val
}

func (c *CLI) AddUintVar(name string, defValue uint, usage string, opts *FlagOptions) *uint {
	val := &defValue
	c.addVar(name, val, strconv.FormatUint(uint64(defValue), 10), usage, opts)
	return val
}

func (c *CLI) AddDurationVar(name, defValue string, usage string, opts *FlagOptions) *time.Duration {
	if opts != nil {
		if opts.PreHook != nil {
			var err error
			defValue, err = opts.PreHook(defValue)
			if err != nil {
				panic("cli: add duration variable \"" + name + "\": " + err.Error())
			}
		}
	}

	valDuration, err := parseDuration(defValue)
	if err != nil {
		panic("cli: add duration variable \"" + name + "\": " + err.Error())
	}

	val := &valDuration
	c.addVar(name, val, defValue, usage, opts)
	return val
}

func writeVar(val string, to interface{}, preHook func(string) (string, error)) error {
	if preHook != nil {
		var err error
		val, err = preHook(val)
		if err != nil {
			return err
		}
	}

	switch to := to.(type) {
	case *string:
		*to = val

	case *int:
		val, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		*to = val

	case *bool:
		val := true
		*to = val

	case *uint:
		val, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		*to = uint(val)

	case *time.Duration:
		val, err := parseDuration(val)
		if err != nil {
			return err
		}
		*to = val

	default:
		panic("cli: write variable: unknown \"to\" argument type")
	}

	return nil
}

func (c *CLI) printVersion() {
	fmt.Println(c.version)
	os.Exit(0)
}

func (c *CLI) printHelp() {
	// Search for the longest flag and required flags list.
	var maxFlagSize int
	var reqFlags string

	for _, v := range c.vars {
		flagSize := len(v.cliFlagName)
		if flagSize > maxFlagSize {
			maxFlagSize = flagSize
		}

		if v.required {
			reqFlags += "[" + v.cliFlagName + "] "
		}
	}

	// Print help
	fmt.Println("Usage:", os.Args[0], reqFlags+"[OPTION]...")
	fmt.Println("")

	for _, v := range c.vars {
		var spaces string
		for i := 0; i < maxFlagSize-len(v.cliFlagName)+2; i++ {
			spaces += " "
		}

		var defaultStr string
		if v.valueDefault != "" {
			defaultStr = " (default: " + v.valueDefault + ")"
		}

		fmt.Println(" ", v.cliFlagName, spaces, v.usage+defaultStr)
	}

	fmt.Println()
	fmt.Println("  -version   Display version and exit.")
	fmt.Println("  -help      Display this help and exit.")

	os.Exit(0)
}

func (c *CLI) Parse() {
	// The name of variables that were read from environment variables or CLI flags.
	// Used to check if "required" flags are present.
	readVars := make(map[string]struct{})

	// Read variables from CLI flags
	{
		alreadyRead := make(map[string]struct{})

		var varInProgress *variable
		for _, arg := range os.Args[1:] {
			if varInProgress == nil {
				switch arg {
				case "-version":
					c.printVersion()

				case "-help":
					c.printHelp()
				}

				_, exist := alreadyRead[arg]
				if exist {
					exitOnError("flag \"" + varInProgress.cliFlagName + "\" occurs twice")
				}

				ok := false
				for _, v := range c.vars {
					if v.cliFlagName == arg {
						switch v.value.(type) {
						case *bool:
							// pass
						default:
							varInProgress = &v
						}

						alreadyRead[v.cliFlagName] = struct{}{}
						readVars[v.name] = struct{}{}

						ok = true
						break
					}
				}

				if !ok {
					exitOnError("unknown flag \"" + arg + "\"")
				}

			} else {
				err := writeVar(arg, varInProgress.value, varInProgress.preHook)
				if err != nil {
					exitOnError("read \"" + varInProgress.cliFlagName + "\" flag: " + err.Error())
				}

				varInProgress = nil
			}
		}

		if varInProgress != nil {
			exitOnError("no value for \"" + varInProgress.cliFlagName + "\" flag")
		}
	}

	// Check required variables
	for _, v := range c.vars {
		if v.required {
			_, ok := readVars[v.name]
			if !ok {
				exitOnError("\"" + v.cliFlagName + "\" flag is missing")
			}
		}
	}
}
