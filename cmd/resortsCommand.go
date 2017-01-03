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

package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/nlamirault/chione/skiinfo"
)

// ResortsCommand is the command which display available resorts
var ResortsCommand = cli.Command{
	Name: "resorts",
	Subcommands: []cli.Command{
		resortsListCommand,
	},
}

var resortsListCommand = cli.Command{
	Name:  "list",
	Usage: "List all resorts",
	Action: func(context *cli.Context) error {

		fmt.Println("Resorts:")
		resorts, err := skiinfo.ListResorts("france")
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return nil
		}
		for _, resort := range resorts {
			fmt.Printf("- %s [%s]\n", resort.Name, resort.Region)
		}
		return nil
	},
}
