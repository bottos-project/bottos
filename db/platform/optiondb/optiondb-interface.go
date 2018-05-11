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
 * file description: option database interface
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */
package optiondb

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OptionDbRepo interface {
	InsertOptionDb(collection string, value interface{}) error
	OptionDbFind(collection string, key string, value interface{}) (interface{}, error)
}

var insertSession *MongoContext
var getSession *MongoContext

func (r *OptionDbRepository) InsertOptionDb(collection string, value interface{}) error {
	var err error
	if insertSession == nil || insertSession.mgoSession == nil {
		insertSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			fmt.Println("collection cccccccccccc", insertSession, collection, err)
			return errors.New("Get session faild" + r.mgoEndpoint)
		}
	}

	insert := func(c *mgo.Collection) error {
		return c.Insert(value)
	}
	insertSession.SetCollection(collection, insert)
	//insertSession.Close()

	return nil
}

func (r *OptionDbRepository) RemoveAllOptionDb(collection string) error {
	var err error
	if getSession == nil || getSession.mgoSession == nil {
		getSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			fmt.Println("collection ", getSession, collection, err)
			return errors.New("Get session faild" + r.mgoEndpoint)
		}
	}
	removeAll := func(c *mgo.Collection) error {
		_, err := c.RemoveAll(nil)
		return err
	}
	getSession.SetCollection(collection, removeAll)

	return nil
}

func (r *OptionDbRepository) OptionDbFind(collection string, key string, value interface{}) (interface{}, error) {
	var err error
	if getSession == nil || getSession.mgoSession == nil {
		getSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			fmt.Println("collection ", getSession, collection, err)
			return nil, errors.New("Get session faild" + r.mgoEndpoint)
		}
	}
	var mesgs interface{}
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{key: value}).All(&mesgs)
	}
	getSession.SetCollection(collection, query)
	return mesgs, nil
}
