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

// ResortCommand is the command which manage a resort
var ResortCommand = cli.Command{
	Name: "resort",
	Subcommands: []cli.Command{
		resortDescribeCommand,
	},
}

var resortDescribeCommand = cli.Command{
	Name:  "describe",
	Usage: "Describe current resort",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "resort",
			Usage: "the resort name",
		},
		cli.StringFlag{
			Name:  "region",
			Usage: "the region name",
		},
	},
	Action: func(context *cli.Context) error {
		if !context.IsSet("resort") {
			return fmt.Errorf("Please specify the resort to used via the --resort option")
		}
		if !context.IsSet("region") {
			return fmt.Errorf("Please specify the region to used via the --region option")
		}
		resort, err := skiinfo.GetResort(context.String("resort"), context.String("region"))
		if err != nil {
			fmt.Println(redOut(err))
			return nil
		}
		fmt.Printf("Resort:\n")
		fmt.Printf("Status: %s\n", resort.Status)
		fmt.Printf("Snow depth piste: low:%s, middle:%s, up:%s\n",
			resort.Piste.Lower, resort.Piste.Middle, resort.Piste.Upper)
		fmt.Printf("Snow depth off-piste: low:%s, middle:%s, up:%s\n",
			resort.OffPiste.Lower, resort.OffPiste.Middle, resort.OffPiste.Upper)
		return nil
	},
}
