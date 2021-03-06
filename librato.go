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

import (
	"fmt"
	"net/url"
	"net/http"
)

const (
	LibratoMetricsPostUrl = "https://%s:%s@metrics-api.librato.com/v1/metrics"
)

type Librato struct {
	Logger
	postMetricsUrl string
}

func NewLibrato(apiUser, apiToken string, logger Logger) *Librato {
	return &Librato{ Logger: logger, postMetricsUrl: assembleLibratoUrl(LibratoMetricsPostUrl, apiUser, apiToken) }
}

func assembleLibratoUrl(pattern, apiUser, apiToken string) string {
	return fmt.Sprintf(pattern, url.QueryEscape(apiUser), apiToken)
}

type libratoMsg struct {
	Gauges []libratoMetric `json:"gauges,omitempty"`
	Counters []libratoMetric `json:"counters,omitempty"`
}

type libratoMetric struct {
	Name string `json:"name"`
	Value float64 `json:"value"`
	Source string `json:"source"`
}

// This method can be used as the Metrics relay function.
func (self *Librato) SendMetricsToLibrato(sourceName string, metrics []Metric) {
	msg := libratoMsg{}

	for i := range metrics {
		metric := libratoMetric{ Name: metrics[i].Name, Value: metrics[i].Value, Source: sourceName }
		switch metrics[i].Type {
			case Counter: msg.Counters = append(msg.Counters, metric)
			case Gauge: msg.Gauges = append(msg.Gauges, metric)
		}
	}

	if statusCode, response, err := HttpPostJson(self.postMetricsUrl, msg, nil); err != nil {
		self.Logf(Warn, "Unable to send metrics to librato - error: %v - response: %s", err, string(response))
	} else if statusCode != http.StatusOK {
		self.Logf(Warn, "Non-200 status code returned by librato for post - response: %s", string(response))
	}
}

