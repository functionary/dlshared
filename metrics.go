/**
 * (C) Copyright 2014, Deft Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dlshared

import "time"

type MetricsRelayFunction func(string, []Metric)

type metricType int8

const (
	Counter metricType = 0
	Gauge metricType = 1

	CounterStr = "counter"
	GaugeStr = "gauge"
)

type Metric struct {
	Name string
	Type metricType
	Value float64
}

type Metrics struct {
	sourceName string
	quitChannel chan bool
	relayFuncs []MetricsRelayFunction
	relayPeriodInSecs int
	metricChannel chan *Metric
	ticker *time.Ticker
}

// The relay function is only called if there are metrics to relay.
func NewMetrics(	sourceName string,
					relayFuncs []MetricsRelayFunction,
					relayPeriodInSecs int,
					metricQueueLength int) *Metrics {

	return &Metrics{
		sourceName: sourceName,
		relayFuncs: relayFuncs,
		relayPeriodInSecs: relayPeriodInSecs,
		quitChannel: make(chan bool),
		metricChannel: make(chan *Metric, metricQueueLength),
	}
}

// Update the gauge value
func (self *Metrics) Gauge(metricName string, value float64) {
	self.metricChannel <- &Metric{
		Name: metricName,
		Type: Gauge,
		Value: value,
	}
}

// Increases the counter by one.
func (self *Metrics) Count(metricName string) {
	self.metricChannel <- &Metric{
		Name: metricName,
		Type: Counter,
		Value: 1,
	}
}

// Increase the counter
func (self *Metrics) CountWithValue(metricName string, value float64) {
	self.metricChannel <- &Metric{
		Name: metricName,
		Type: Counter,
		Value: value,
	}
}

func (self *Metrics) listenForEvents() {

	metrics := make(map[string]*Metric)

    for {
        select {
			case metric := <- self.metricChannel:
				current, found := metrics[metric.Name]
				if !found { metrics[metric.Name] = metric; continue }

				if metric.Type == Counter { current.Value = current.Value + metric.Value
				} else {  current.Value = metric.Value }

			case <- self.ticker.C:

				var toRelay []Metric

				for _, v := range metrics { toRelay = append(toRelay, Metric{ Name: v.Name, Type: v.Type, Value: v.Value }) }

				if len(toRelay) == 0 { continue }

				for _, relayFunc := range self.relayFuncs { go relayFunc(self.sourceName, toRelay) }

			case <- self.quitChannel:
				self.ticker.Stop()
				return
        }
    }
}


func (self *Metrics) Start() error {

	self.ticker = time.NewTicker(time.Duration(self.relayPeriodInSecs) * time.Second)

	go self.listenForEvents()

	return nil
}

func (self *Metrics) Stop() error {
	self.quitChannel <- true
	return nil
}

// Pass the logger struct to the log metrics.
func NewLogMetrics(logger Logger) *LogMetrics {
	return &LogMetrics{ Logger: logger, enabled: true }
}

// This can be used to print metrics out to a configured logger.
type LogMetrics struct {
	Logger
	enabled bool
}

func (self *LogMetrics) Disable() { self.enabled = false }

func (self *LogMetrics) Enable() { self.enabled = true }

// This logs an info message with the following format: [source: %s - type: %s - metric: %s - value: %f]
func (self *LogMetrics) Log(sourceName string, metrics []Metric) {

	if !self.enabled { return }

	var typeStr string
	for i := range metrics {
		switch metrics[i].Type {
			case Counter: typeStr = CounterStr
			case Gauge: typeStr = GaugeStr
		}
		self.Logf(Info, "[source: %s - type: %s - metric: %s - value: %f]", sourceName, typeStr, metrics[i].Name, metrics[i].Value)
	}
}

