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
	"labix.org/v2/mgo/bson"
)

// A metrics relay function that stores count and guage values in mongo. This
// does not store historical values, simply current ones.
type MongoMetrics struct {
	Logger
	DataSource
	mongoComponentName string
	fireAndForget bool
}

/*
type persistedMongoMetric struct {
	Id string `bson:"_id"`

	Source string `bson:"source"`
	Name string `bson:"name"`
	Type string `bson:"type"`
	Value float64 `bson:"value"`

	Updated *time.Time `bson:"updated"` // The last update time
	Created *time.Time `bson:"created"`
}
*/

func NewMongoMetrics(dbName, collectionName, mongoComponentName string, fireAndForget bool) *MongoMetrics {
	return &MongoMetrics{ Logger: Logger{}, DataSource: DataSource{ DbName: dbName, CollectionName: collectionName }, mongoComponentName: mongoComponentName, fireAndForget: fireAndForget }
}

// Assemble the doc id. If there is an error, it is logged here.
func (self *MongoMetrics) assembleDocId(metricName, sourceName string) (string, error) {
	id, err := Md5Hex(fmt.Sprintf("%s-%s-metrics", metricName, sourceName))
	if err != nil {
		self.Logf(Error, "Unable to assemble doc id - metric: %s - source: %s - error: %v", metricName, sourceName, err)
		return nadaStr, err
	}

	return id, nil
}

func (self *MongoMetrics) persistCounter(sourceName string, metric *Metric) {

	docId, err := self.assembleDocId(metric.Name, sourceName)
	if err != nil { return }

	selector := &bson.M{ "_id": docId }

	now := self.Now()

	upsert := &bson.M{
		"$setOnInsert": &bson.M{ "name": metric.Name, "source": sourceName, "type": CounterStr, "created": now, "updated": now },
		"$set": &bson.M{ "updated": now },
		"$inc": &bson.M{ "value": metric.Value },
	}

	if self.fireAndForget { err = self.Upsert(selector, upsert)
	} else { err = self.UpsertSafe(selector, upsert) }

	if err != nil { self.Logf(Error, "Unable to persist counter - source: %s - metric: %s - error: %v", sourceName, metric.Name, err) }
}

func (self *MongoMetrics) persistGauge(sourceName string, metric *Metric) {

	docId, err := self.assembleDocId(metric.Name, sourceName)
	if err != nil { return }

	selector := &bson.M{ "_id": docId }

	now := self.Now()

	upsert := &bson.M{
		"$setOnInsert": &bson.M{ "name": metric.Name, "source": sourceName, "type": CounterStr, "created": now, "updated": now },
		"$set": &bson.M{ "updated": now, "value": metric.Value },
	}

	if self.fireAndForget { err = self.Upsert(selector, upsert)
	} else { err = self.UpsertSafe(selector, upsert) }

	if err != nil { self.Logf(Error, "Unable to persist counter - source: %s - metric: %s - error: %v", sourceName, metric.Name, err) }
}

// This method can be used as the Metrics relay function.
func (self *MongoMetrics) StoreMetricsInMongo(sourceName string, metrics []Metric) {

	for i := range metrics {
		switch metrics[i].Type {
			case Counter: self.persistCounter(sourceName, &metrics[i])
			case Gauge: self.persistGauge(sourceName, &metrics[i])
		}
	}
}

func (self *MongoMetrics) Start(kernel *Kernel) error {

	self.Logger = kernel.Logger
	self.Mongo = kernel.GetComponent(self.mongoComponentName).(*Mongo)

	if err := self.EnsureIndex([]string{ "name" }); err != nil { return err }

	if err := self.EnsureIndex([]string{ "source" }); err != nil { return err }

	if err := self.EnsureIndex([]string{ "name", "source", "value" }); err != nil { return err }

	if err := self.EnsureIndex([]string{ "source", "updated" }); err != nil { return err }

	if err := self.EnsureUniqueIndex([]string{ "name", "source" }); err != nil { return err }

	return nil
}

func (self *MongoMetrics) Stop(kernel *Kernel) error { return nil }

