// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description: database option for dapp, store all the data that is conformed
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package optiondb

import (
	"encoding/json"
	"errors"
	//"fmt"

	"github.com/bottos-project/bottos/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	log "github.com/cihub/seelog"
)

//OptionDbRepository is the option db struct
type OptionDbRepository struct {
	mgoEndpoint string
	is_optiondb_offline bool
}

//NewOptionDbRepository creates a new OptionDbRepository
func NewOptionDbRepository(endpoint string) *OptionDbRepository {
	return &OptionDbRepository{mgoEndpoint: endpoint, is_optiondb_offline: false}
}

//MongoContext is a plugin db for option db, you can chose different db if you like.
type MongoContext struct {
	mgoSession *mgo.Session
}

//GetSession is to create session to mongodb
func GetSession(url string) (*MongoContext, error) {
	if url == "" {
		log.Error("Error! GetSession failed ! Url is empty")
		return nil, errors.New("invalid para url")
	}
	var err error
	// tried doing this - doesn't work as intended
	defer func() {
		if r := recover(); r != nil {
			log.Error("Detected panic")
			var ok bool
			err, ok := r.(error)
			if !ok {
				log.Errorf("pkg:  %v,  error: %s", r, err)
			}
		}
	}()

	//maxWait := time.Duration(5 * time.Second)
	mgoSession, err := mgo.Dial(url)
	if err != nil {
		log.Error(err)
		return nil, errors.New("Dial faild" + url)
	}
	return &MongoContext{mgoSession.Clone()}, nil
}

//GetCollection is to get mongodb collection
func (c *MongoContext) GetCollection(db string, collection string) *mgo.Collection {
	session := c.mgoSession
	//defer session.Close()
	collects := session.DB(config.DEFAULT_OPTIONDB_NAME).C(collection)
	return collects
}

//SetCollection is to set mongodb collection
func (c *MongoContext) SetCollection(collection string, s func(*mgo.Collection) error) error {
	session := c.mgoSession
	//defer session.Close()
	collects := session.DB(config.DEFAULT_OPTIONDB_NAME).C(collection)
	return s(collects)
}

//SetCollectionCount is to set mongodb collection by returning records number
func (c *MongoContext) SetCollectionCount(collection string, s func(*mgo.Collection) (int, error)) (int, error) {
	session := c.mgoSession
	defer session.Close()
	collects := session.DB(config.DEFAULT_OPTIONDB_NAME).C(collection)
	return s(collects)
}

//SetCollectionByDB is to set mongodb collection by specific database
func (c *MongoContext) SetCollectionByDB(db string, collection string, s func(*mgo.Collection) error) error {
	session := c.mgoSession
	defer session.Close()
	collects := session.DB(db).C(collection)
	return s(collects)
}

//Close is to close mongodb session
func (c *MongoContext) Close() {
	c.mgoSession.Close()
}

// CollectionExists returns true if the collection name exists in the specified database.
func (c *MongoContext) isCollectionExists(useCollection string) bool {
	session := c.mgoSession
	database := session.DB(config.DEFAULT_OPTIONDB_NAME)
	collections, err := database.CollectionNames()
	if err != nil {
		return false
	}

	for _, collection := range collections {
		if collection == useCollection {
			return true
		}
	}

	return false
}

// ToString converts the quer map to a string.
func ToString(queryMap interface{}) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// ToStringD converts bson.D to a string.
func ToStringD(queryMap bson.D) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}
