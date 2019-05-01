/*
 *    Copyright (C) 2019 CalmBit
 *                  2014-2019 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      CalmBit <calmbit@posteo.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package prometheusbee is a bee designed to expose metrics for scraping
// by Prometheus.
package prometheusbee

import (
	"github.com/muesli/beehive/bees"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// PrometheusBee is a bee designed to expose metrics for scraping
// by Prometheus.
type PrometheusBee struct {
	bees.Bee

	addr string

	counterVecName   string
	gaugeVecName     string
	histogramVecName string
	summaryVecName   string

	counter   *prometheus.CounterVec
	gauge     *prometheus.GaugeVec
	histogram *prometheus.HistogramVec
	summary   *prometheus.SummaryVec
}

// Run executes the Bee's event loop.
func (mod *PrometheusBee) Run(eventChan chan bees.Event) {

	// Counter vector registration

	mod.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: mod.counterVecName,
			Help: "Collection of all counter variables being iterated by hives",
		},
		[]string{"name"},
	)
	err := prometheus.Register(mod.counter)
	if err != nil {
		mod.LogErrorf("Error registering counter vector: %v", err)
		panic("Unable to start Prometheus due to counter initialization failure")
	}

	// Gauge vector registration

	mod.gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: mod.gaugeVecName,
			Help: "Collection of all gauge variables being iterated by hives",
		},
		[]string{"name"},
	)
	err = prometheus.Register(mod.gauge)
	if err != nil {
		mod.LogErrorf("Error registering gauge vector: %v", err)
		panic("Unable to start Prometheus due to gauge initialization failure")
	}

	// Histogram vector registration

	mod.histogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: mod.histogramVecName,
			Help: "Collection of all histogram variables being iterated by hives",
		},
		[]string{"name"},
	)
	err = prometheus.Register(mod.histogram)
	if err != nil {
		mod.LogErrorf("Error registering histogram vector: %v", err)
		panic("Unable to start Prometheus due to histogram initialization failure")
	}

	// Summary vector registration

	mod.summary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: mod.summaryVecName,
			Help: "Collection of all summary variables being iterated by hives",
		},
		[]string{"name"},
	)
	err = prometheus.Register(mod.summary)
	if err != nil {
		mod.LogErrorf("Error registering summary vector: %v", err)
		panic("Unable to start Prometheus due to summary initialization failure")
	}

	// Now, to serve everything up:

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			mod.LogErrorf("Error running prometheus metric handler: %v", err)
		}
	}()

	select {
	case <-mod.SigChan:
		return
	}
}

// Action triggers the action passed to it.
func (mod *PrometheusBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "counter_inc":
		var label string
		action.Options.Bind("label", &label)
		c, err := mod.counter.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error incrementing counter: %v", err)
		} else {
			c.Inc()
		}
	case "counter_add":
		var label string
		var value float64
		action.Options.Bind("label", &label)
		action.Options.Bind("value", &value)
		c, err := mod.counter.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error adding to counter: %v", err)
		} else {
			c.Add(value)
		}
	case "gauge_set":
		var label string
		var value float64
		action.Options.Bind("label", &label)
		action.Options.Bind("value", &value)
		g, err := mod.gauge.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error setting gauge: %v", err)
		} else {
			g.Set(value)
		}
	case "gauge_inc":
		var label string
		action.Options.Bind("label", &label)
		g, err := mod.gauge.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error incrementing gauge: %v", err)
		} else {
			g.Inc()
		}
	case "gauge_dec":
		var label string
		action.Options.Bind("label", &label)
		g, err := mod.gauge.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error decrementing gauge: %v", err)
		} else {
			g.Dec()
		}
	case "gauge_add":
		var label string
		var value float64
		action.Options.Bind("label", &label)
		action.Options.Bind("value", &value)
		g, err := mod.gauge.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error adding to gauge: %v", err)
		} else {
			g.Add(value)
		}
	case "gauge_sub":
		var label string
		var value float64
		action.Options.Bind("label", &label)
		action.Options.Bind("value", &value)
		g, err := mod.gauge.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error subtracting from gauge: %v", err)
		} else {
			g.Sub(value)
		}
	case "gauge_set_to_current_time":
		var label string
		action.Options.Bind("label", &label)
		g, err := mod.gauge.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error setting gauge to current time: %v", err)
		} else {
			g.SetToCurrentTime()
		}
	case "histogram_observe":
		var label string
		var value float64
		action.Options.Bind("label", &label)
		action.Options.Bind("value", &value)
		h, err := mod.histogram.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error observing value for histogram: %v", err)
		} else {
			h.Observe(value)
		}
	case "summary_observe":
		var label string
		var value float64
		action.Options.Bind("label", &label)
		action.Options.Bind("value", &value)
		s, err := mod.summary.GetMetricWithLabelValues(label)
		if err != nil {
			mod.LogErrorf("Error observing value for summary: %v", err)
		} else {
			s.Observe(value)
		}
	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *PrometheusBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.addr)
	options.Bind("counter_vec_name", &mod.counterVecName)
	options.Bind("gauge_vec_name", &mod.gaugeVecName)
	options.Bind("histogram_vec_name", &mod.histogramVecName)
	options.Bind("summary_vec_name", &mod.summaryVecName)
}
