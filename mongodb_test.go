/*
 * Copyright 2016 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eval

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

const (
	colName = "test"
)

type Foo struct {
	ID      bson.ObjectId `bson:"_id"`
	StrVal  string
	IntVal  int
	BoolVal bool
}

func TestMongoDBEnvironmentSession(t *testing.T) {
	testDoc := Foo{
		bson.NewObjectId(),
		"Ipsum",
		999,
		true,
	}

	env := PrepareMongoDBEnvironment(t)
	if env == nil {
		return
	}
	defer env.Dispose()

	session := env.Session()

	col := session.DB("").C(colName)
	if err := col.Insert(testDoc); err != nil {
		t.Fatalf("Error inserting new document: %v", err)
	}

	var getDoc Foo
	if err := col.FindId(testDoc.ID).One(&getDoc); err != nil {
		t.Fatalf("Error getting inserted document: %v", err)
	}

	if getDoc.ID != testDoc.ID ||
		getDoc.StrVal != testDoc.StrVal ||
		getDoc.IntVal != testDoc.IntVal ||
		getDoc.BoolVal != testDoc.BoolVal {
		t.Error("Inserted document does match the original one")
		t.Errorf("Original: %#v", testDoc)
		t.Errorf("Inserted: %#v", getDoc)
		t.FailNow()
	}
}
