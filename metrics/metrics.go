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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "chione"
)

var (
	lower = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "lower"),
		"Snow depths at the lower of the ski resort",
		nil, nil,
	)
	middle = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "middle"),
		"Snow depths in the middle of the ski resort",
		nil, nil,
	)
	upper = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upper"),
		"Snow depths at the upper of the ski resort).",
		nil, nil,
	)
)

// Exporter collects metrics from the given ski resort and exports them using
// the prometheus metrics package.
type Exporter struct {
	SkiResortName   string
	SkiResortRegion string
}

// NewExporter returns an initialized Exporter.
func NewExporter(name string, region string) (*Exporter, error) {
	log.Debugln("Init exporter")
	return &Exporter{
		SkiResortName:   name,
		SkiResortRegion: region,
	}, nil
}

// Describe describes all the metrics ever exported by the Speedtest exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- lower
	ch <- middle
	ch <- upper
}

// Collect fetches the stats from configured Speedtest location and delivers them
// as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Infof("Exporter starting")

	log.Infof("Speedtest exporter finished")
}
