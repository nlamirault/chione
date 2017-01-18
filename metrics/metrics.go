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

	if len(resort.Piste.Lower) > 0 {
		lowerVal, err := strconv.ParseFloat(strings.Replace(resort.Piste.Lower, "cm", "", -1), 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %d", e.SkiResortRegion, e.SkiResortName, resort.Piste.Lower)
			return
		}
		ch <- prometheus.MustNewConstMetric(lower, prometheus.GaugeValue, lowerVal, e.SkiResortName)
	}

	if len(resort.Piste.Middle) > 0 {
		middleVal, err := strconv.ParseFloat(strings.Replace(resort.Piste.Middle, "cm", "", -1), 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %d", e.SkiResortRegion, e.SkiResortName, resort.Piste.Middle)
			return
		}
		ch <- prometheus.MustNewConstMetric(middle, prometheus.GaugeValue, middleVal, e.SkiResortName)
	}

	if len(resort.Piste.Upper) > 0 {
		upperVal, err := strconv.ParseFloat(strings.Replace(resort.Piste.Upper, "cm", "", -1), 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %d", e.SkiResortRegion, e.SkiResortName, resort.Piste.Upper)
			return
		}
		ch <- prometheus.MustNewConstMetric(upper, prometheus.GaugeValue, upperVal, e.SkiResortName)
	}
}

func (e *Exporter) collectSlopes(ch chan<- prometheus.Metric, resort *skiinfo.ResortDescription) {
	beginnerSlopes := resort.Slopes.Beginning.String()
	if len(beginnerSlopes) > 0 {
		tokens := strings.Split(beginnerSlopes, "/")
		val, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %s %s", e.SkiResortRegion, e.SkiResortName, beginnerSlopes, tokens)
			return
		}
		ch <- prometheus.MustNewConstMetric(beginner, prometheus.GaugeValue, val, e.SkiResortName)
	}
	intermediateSlopes := resort.Slopes.Intermediate.String()
	if len(intermediateSlopes) > 0 {
		tokens := strings.Split(intermediateSlopes, "/")
		val, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %s %s", e.SkiResortRegion, e.SkiResortName, intermediateSlopes, tokens)
			return
		}
		ch <- prometheus.MustNewConstMetric(intermediate, prometheus.GaugeValue, val, e.SkiResortName)
	}
	advancedSlopes := resort.Slopes.Advanced.String()
	if len(advancedSlopes) > 0 {
		tokens := strings.Split(advancedSlopes, "/")
		val, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %s %s", e.SkiResortRegion, e.SkiResortName, advancedSlopes, tokens)
			return
		}
		ch <- prometheus.MustNewConstMetric(advanced, prometheus.GaugeValue, val, e.SkiResortName)
	}
	exporterSlopes := resort.Slopes.Expert.String()
	if len(exporterSlopes) > 0 {
		tokens := strings.Split(exporterSlopes, "/")
		val, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			log.Errorf("Can't parse value : %s %s %s %s", e.SkiResortRegion, e.SkiResortName, exporterSlopes, tokens)
			return
		}
		ch <- prometheus.MustNewConstMetric(expert, prometheus.GaugeValue, val, e.SkiResortName)
	}

}
