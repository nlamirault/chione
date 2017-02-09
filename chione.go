// Copyright (C) 2017 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	//"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/nlamirault/chione/cmd"
	"github.com/nlamirault/chione/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "chione"
	app.Usage = "CLI for skiing resorts informations"
	app.Version = version.Version
	// app.Copyright = "(c) 2017 Nicolas Lamirault"
	app.Commands = []cli.Command{
		cmd.VersionCommand,
		cmd.ResortsCommand,
		cmd.ResortCommand,
		cmd.MetricsCommand,
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
			// Value: false,
			Usage: "Enable debug mode",
		},
	}
	app.Action = func(context *cli.Context) error {
		cli.ShowAppHelp(context)
		return nil
	}
	app.Before = func(context *cli.Context) error {
		logrus.SetOutput(os.Stdout)
		if context.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.WarnLevel)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
