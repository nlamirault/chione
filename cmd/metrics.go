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
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"

	"github.com/nlamirault/chione/config"
	"github.com/nlamirault/chione/metrics"
)

// MetricsCommand is the command which manage metrics for ski resorts
var MetricsCommand = cli.Command{
	Name: "metrics",
	Subcommands: []cli.Command{
		metricsDryRunCommand,
		metricsExportCommand,
	},
}

var metricsDryRunCommand = cli.Command{
	Name:  "dryRun",
	Usage: "Test metrics for a ski resort",
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
		displaySkiResortMetrics(context.String("resort"), context.String("region"))
		return nil
	},
}
var metricsExportCommand = cli.Command{
	Name:  "export",
	Usage: "Export metrics for a ski resort",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "the configuration file to use",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "Port to listen on for web interface and telemetry.",
			Value: 9114,
		},
		cli.StringFlag{
			Name:  "metricsPath",
			Usage: "Path under which to expose metrics",
			Value: "/metrics",
		},
	},
	Action: func(context *cli.Context) error {
		if !context.IsSet("config") {
			return fmt.Errorf("Please specify the configuration file to used via the --config option")
		}
		exportSkiResortMetrics(context.String("config"), context.Int("port"), context.String("metricsPath"))
		return nil
	},
}

func displaySkiResortMetrics(name string, region string) {
}

func exportSkiResortMetrics(configFile string, port int, metricsPath string) {
	conf, err := config.New(configFile)
	if err != nil {
		fmt.Println(redOut(err))
		return
	}
	exporter, err := metrics.NewExporter(conf)
	if err != nil {
		fmt.Println(redOut(err))
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Chione Exporter</title></head>
             <body>
             <h1>Chione Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	fmt.Println("Listening on", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
