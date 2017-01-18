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
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/nlamirault/chione/skiinfo"
)

const (
	namespace = "chione"
)

var (

	// Snow depths

	lower = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "lower"),
		"Snow depths at the lower of the ski resort",
		[]string{"name"}, nil,
	)
	middle = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "middle"),
		"Snow depths in the middle of the ski resort",
		[]string{"name"}, nil,
	)
	upper = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upper"),
		"Snow depths at the upper of the ski resort.",
		[]string{"name"}, nil,
	)

	// Slopes

	beginner = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "beginner"),
		"Number of open beginner slopes.",
		[]string{"name"}, nil,
	)
	intermediate = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "intermediate"),
		"Number of open intermediate slopes.",
		[]string{"name"}, nil,
	)
	advanced = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "advanced"),
		"Number of open advanced slopes.",
		[]string{"name"}, nil,
	)
	expert = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "expert"),
		"Number of open expert slopes.",
		[]string{"name"}, nil,
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
	resort, err := skiinfo.GetResort(e.SkiResortName, e.SkiResortRegion)
	if err != nil {
		log.Errorf("Can't retrive metrics: %s", err.Error())
		return
	}
	e.collectSnowDepth(ch, resort)
	e.collectSlopes(ch, resort)
	log.Infof("Exporter finished")
}

func (e *Exporter) collectSnowDepth(ch chan<- prometheus.Metric, resort *skiinfo.ResortDescription) {
	log.Infof("Ski resort informations: %s", resort.Piste)
	e.addSnowDepth(ch, lower, resort.Piste.Lower)
	e.addSnowDepth(ch, middle, resort.Piste.Middle)
	e.addSnowDepth(ch, upper, resort.Piste.Upper)
}

func (e *Exporter) addSnowDepth(ch chan<- prometheus.Metric, desc *prometheus.Desc, value string) {
	if len(value) > 0 {
		val, err := strconv.ParseFloat(strings.Replace(value, "cm", "", -1), 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %s", e.SkiResortRegion, e.SkiResortName, value)
			return
		}
		log.Debugf("Add snow depth metric %d to desc %s", val, desc)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, e.SkiResortName)
	}

}

func (e *Exporter) collectSlopes(ch chan<- prometheus.Metric, resort *skiinfo.ResortDescription) {
	log.Infof("Ski resort slopes: %s", resort.Slopes)
	e.addSlopes(ch, beginner, resort.Slopes.Beginning.String())
	e.addSlopes(ch, intermediate, resort.Slopes.Intermediate.String())
	e.addSlopes(ch, advanced, resort.Slopes.Advanced.String())
	e.addSlopes(ch, expert, resort.Slopes.Expert.String())
}

func (e *Exporter) addSlopes(ch chan<- prometheus.Metric, desc *prometheus.Desc, value string) {
	if len(value) > 0 {
		tokens := strings.Split(value, "/")
		val, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %s %s", e.SkiResortRegion, e.SkiResortName, value, tokens)
			return
		}
		log.Debugf("Add slopes metric %d to desc %s", val, desc)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, e.SkiResortName)
	}
}
