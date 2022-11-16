package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const License = `
Copyright (C) 2021 NETWAYS GmbH <info@netways.de>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
`

func main() {
	config, err := ParseConfigFromFlags(os.Args[1:])
	if err != nil {
		if errors.Is(err, ErrVersionRequested) || errors.Is(err, flag.ErrHelp) {
			os.Exit(check.Unknown)
		}

		check.ExitError(err)
	}

	if config.Debug {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}

	api := RestAPI{URL: config.API, Client: config.NewClient()}

	result, err := api.ExecuteCheck(config.Command, config.Arguments, config.Timeout)
	if err != nil {
		check.ExitError(err)
	}

	_, _ = fmt.Fprintln(os.Stdout, result.String())

	os.Exit(result.ExitCode)
}
