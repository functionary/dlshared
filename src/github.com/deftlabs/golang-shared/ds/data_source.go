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

package deftlabsds

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"github.com/deftlabs/golang-shared/log"
)

type DataSource struct {
	DbName string
	CollectionName string
	Mongo *Mongo
	slogger.Logger
}

// Insert a document into a collection with the base configured write concern.
func (self *DataSource) Insert(doc interface{}) error { return self.Mongo.Collection(self.DbName, self.CollectionName).Insert(doc) }

// Upsert a document in a collection with the base configured write concern.
func (self *DataSource) Upsert(selector interface{}, change interface{}) error {
	_, err := self.Mongo.Collection(self.DbName, self.CollectionName).Upsert(selector, change)
	return err
}

// Insert a document into a collection with the passed write concern.
func (self *DataSource) InsertSafe(doc interface{}, safeMode *mgo.Safe) error {
	session := self.SessionClone()

	defer session.Close()

	session.SetSafe(safeMode)
	return session.DB(self.DbName).C(self.CollectionName).Insert(doc)
}

// Finds one document or returns nil. If the document is not found the two return values are nil. This
// method takes the result param and returns it if found (nil otherwise).
func (self *DataSource) FindOne(query *bson.M, result interface{}) (interface{}, error) {
	err := self.Collection().Find(query).One(result)

	if err != nil && err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete one or more documents from the collection. If the document(s) is/are not found, no error
// is returned.
func (self *DataSource) Delete(selector interface{}) error {
	_, err := self.Collection().RemoveAll(selector)
	return err
}

// Ensure a unique, non-sparse index is created. This does not create in the background. This does
// NOT drop duplicates if they exist. Duplicates will cause an error.
func (self *DataSource) EnsureUniqueIndex(fields []string) error {
	return self.Collection().EnsureIndex(mgo.Index{
		Key: fields,
		Unique: true,
		DropDups: true,
		Background: false,
		Sparse: false,
	})
}

// Ensure a non-unique, non-sparse index is created. This does not create in the background.
func (self *DataSource) EnsureIndex(fields []string) error {
	return self.Collection().EnsureIndex(mgo.Index{
		Key: fields,
		Unique: false,
		DropDups: true,
		Background: false,
		Sparse: false,
	})
}

// Ensure a non-unique, sparse index is created. This does not create in the background.
func (self *DataSource) EnsureSparseIndex(fields []string) error {
	return self.Collection().EnsureIndex(mgo.Index{
		Key: fields,
		Unique: false,
		DropDups: true,
		Background: false,
		Sparse: true,
	})
}

// Ensure a unique, sparse index is created. This does not create in the background. This does
// NOT drop duplicates if they exist. Duplicates will cause an error.
func (self *DataSource) EnsureUniqueSparseIndex(fields []string) error {
	return self.Collection().EnsureIndex(mgo.Index{
		Key: fields,
		Unique: true,
		DropDups: false,
		Background: false,
		Sparse: true,
	})
}

// Returns the collection from the session.
func (self *DataSource) Collection() *mgo.Collection { return self.Mongo.Collection(self.DbName, self.CollectionName) }

// Returns the database from the session.
func (self *DataSource) Db() *mgo.Database { return self.Mongo.Db(self.DbName) }

// Returns the session struct.
func (self *DataSource) Session() *mgo.Session { return self.Mongo.session }

// Returns a clone of the session struct.
func (self *DataSource) SessionClone() *mgo.Session { return self.Mongo.session.Clone() }

// Returns a copy of the session struct.
func (self *DataSource) SessionCopy() *mgo.Session { return self.Mongo.session.Clone() }
