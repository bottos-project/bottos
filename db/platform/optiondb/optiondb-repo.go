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
	"fmt"
	//"time"

	"github.com/bottos-project/core/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OptionDbRepository struct {
	mgoEndpoint string
}

//NewOptionDbRepository creates a new OptionDbRepository
func NewOptionDbRepository(endpoint string) *OptionDbRepository {
	return &OptionDbRepository{mgoEndpoint: endpoint}
}

type MongoContext struct {
	mgoSession *mgo.Session
}

func GetSession(url string) (*MongoContext, error) {
	if url == "" {
		return nil, errors.New("invalid para url")
	}
	var err error
	// tried doing this - doesn't work as intended
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Detected panic")
			var ok bool
			err, ok := r.(error)
			if !ok {
				fmt.Printf("pkg:  %v,  error: %s", r, err)
			}
		}
	}()

	//maxWait := time.Duration(5 * time.Second)
	mgoSession, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Dial faild" + url)
	}
	return &MongoContext{mgoSession.Clone()}, nil
}
func (c *MongoContext) GetCollection(db string, collection string) *mgo.Collection {
	session := c.mgoSession
	//defer session.Close()
	collects := session.DB(config.DEFAULT_OPTIONDB_NAME).C(collection)
	return collects
}
func (c *MongoContext) SetCollection(collection string, s func(*mgo.Collection) error) error {
	session := c.mgoSession
	//defer session.Close()
	collects := session.DB(config.DEFAULT_OPTIONDB_NAME).C(collection)
	return s(collects)
}

func (c *MongoContext) SetCollectionCount(collection string, s func(*mgo.Collection) (int, error)) (int, error) {
	session := c.mgoSession
	defer session.Close()
	collects := session.DB(config.DEFAULT_OPTIONDB_NAME).C(collection)
	return s(collects)
}
func (c *MongoContext) SetCollectionByDB(db string, collection string, s func(*mgo.Collection) error) error {
	session := c.mgoSession
	defer session.Close()
	collects := session.DB(db).C(collection)
	return s(collects)
}

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
