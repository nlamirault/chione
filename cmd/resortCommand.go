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
	"bytes"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
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
		describeSkiResort(context.String("resort"), context.String("region"))
		return nil
	},
}

func describeSkiResort(name string, region string) {
	resort, err := skiinfo.GetResort(name, region)
	if err != nil {
		fmt.Println(redOut(err))
		return
	}
	// fmt.Printf("Resort:\n")
	table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"", "Club 1", "Club 2", "Score", "Commentaire"})
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	content := []string{"", ""}
	content[0] = yellowOut("Status")
	content[1] = resort.Status
	table.Append(content)

	content = []string{"", ""}
	content[0] = yellowOut("Enneigement sur les pistes") //"Snow depth Piste")
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("- bas: %s\n", resort.Piste.Lower))
	buffer.WriteString(fmt.Sprintf("- milieu: %s\n", resort.Piste.Middle))
	buffer.WriteString(fmt.Sprintf("- haut: %s\n", resort.Piste.Upper))
	content[1] = buffer.String()
	table.Append(content)

	content = []string{"", ""}
	content[0] = yellowOut("Enneigement hors-pistes") // "Snow depth Off-Piste")
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("- bas: %s\n", resort.OffPiste.Lower))
	buffer.WriteString(fmt.Sprintf("- milieu: %s\n", resort.OffPiste.Middle))
	buffer.WriteString(fmt.Sprintf("- haut: %s\n", resort.OffPiste.Upper))
	content[1] = buffer.String()
	table.Append(content)

	content = []string{"", ""}
	content[0] = yellowOut("Chute de neige") //"Snowfall")
	buffer.Reset()
	for k, v := range resort.SnowFall {
		if len(k) == 0 {
			k = "aujourd'hui"
		}
		buffer.WriteString(fmt.Sprintf("- %s: %s\n", k, v))
	}
	content[1] = buffer.String()
	table.Append(content)

	content = []string{"", ""}
	content[0] = yellowOut("Domaine skiable") //"Terrain")
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("- Vertes: %s\n", resort.Slopes.Beginning.String()))
	buffer.WriteString(fmt.Sprintf("- Bleues: %s\n", resort.Slopes.Intermediate.String()))
	buffer.WriteString(fmt.Sprintf("- Rouges: %s\n", resort.Slopes.Advanced.String()))
	buffer.WriteString(fmt.Sprintf("- Noires: %s\n", resort.Slopes.Expert.String()))
	content[1] = buffer.String()
	table.Append(content)

	table.Render()
}
