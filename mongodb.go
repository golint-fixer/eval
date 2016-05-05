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
func PrepareMongoDBEnvironment() (*MongoDBEnvironment, *ErrUser) {
	mongo := NewMongoDBEnvironment()
	if ok, err := mongo.Applicability(); !ok {
		return nil, err
	}

	if ok, err := mongo.Run(); !ok {
		return nil, err
	}

	net, err := mongo.Network()
	if err != nil {
		mongo.Stop()
		return nil, &ErrUser{
			Fatal,
			fmt.Sprintf("Error getting MongoDB IP address: %s", err),
		}
	}

	mgourl := fmt.Sprintf(mongoURLTpl, net[0].IPAddress, net[0].Port)

	session, err := newDBSession(mgourl)
	if err != nil {
		mongo.Stop()
		return nil, &ErrUser{
			Fatal,
			fmt.Sprintf("Error opening a MongoDB session: %s", err),
		}
	}

	return &MongoDBEnvironment{
		mongo,
		session,
	}, nil
}

// Dispose closes database session and removes current environment.
func (e *MongoDBEnvironment) Dispose() {
	if e.session != nil {
		e.session.Close()
		e.session = nil
	}
	if e.env != nil {
		e.env.Stop()
		e.env = nil
	}
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
