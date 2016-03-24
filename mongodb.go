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
	"fmt"
	"testing"

	"gopkg.in/mgo.v2"
)

const (
	mongoURLTpl = "mongodb://%s:%d/auth"
)

// A MongoDBEnvironment represents a MongoDB environment for testing purposes.
type MongoDBEnvironment struct {
	env     *Environment
	session *mgo.Session
}

// PrepareMongoDBEnvironment creates a new MongoDB container, starts it and open
// a new session to database.
func PrepareMongoDBEnvironment(tb testing.TB) *MongoDBEnvironment {
	mongo := NewMongoDBEnvironment(tb)
	if !mongo.Applicability() {
		tb.Skip("This test cannot be run because Docker is not acessible")
		return nil
	}

	if !mongo.Run() {
		tb.Fatal("Could not start MongoDB server")
		return nil
	}

	net, err := mongo.Network()
	if err != nil {
		mongo.Stop()
		tb.Fatalf("Error getting MongoDB IP address: %s\n", err)
		return nil
	}

	mgourl := fmt.Sprintf(mongoURLTpl, net[0].IPAddress, net[0].Port)

	session, err := newDBSession(mgourl)
	if err != nil {
		mongo.Stop()
		tb.Fatalf("Error opening a MongoDB session: %s\n", err)
		return nil
	}

	return &MongoDBEnvironment{
		mongo,
		session,
	}
}

// Dispose closes database session and removes current environment.
func (e *MongoDBEnvironment) Dispose() {
	e.session.Close()
	e.env.Stop()
	e.session = nil
	e.env = nil
}

// Session returns the database session for current environment.
func (e *MongoDBEnvironment) Session() *mgo.Session {
	return e.session
}

func newDBSession(url string) (*mgo.Session, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	_, err = session.DB("").CollectionNames()
	if err != nil {
		return nil, err
	}

	return session, nil
}
