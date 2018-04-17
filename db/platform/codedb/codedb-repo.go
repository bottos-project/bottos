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
 * file description: code database interface
 * @Author: May Luo
 * @Date:   2017-12-05
 * @Last Modified by:
 * @Last Modified time:
 */

package codedb

import (
	"github.com/tidwall/buntdb"
)

type CodeDbRepository struct {
	fn string     // filename for reporting
	db *buntdb.DB // LevelDB instance
}

func NewCodeDbRepository(file string) (*CodeDbRepository, error) {
	codedb, err := buntdb.Open("file")
	if err != nil {
		return nil, err
	}
	return &CodeDbRepository{
		fn: file,
		db: codedb,
	}, nil
}

func (k *CodeDbRepository) CallCreatObject(objectName string, objectValue interface{}) error {
	return nil
}
func (k *CodeDbRepository) CallCreatObjectIndex(objectName string, indexName string) error {
	return nil
}
func (k *CodeDbRepository) CallSetObject(objectName string, objectValue interface{}) error {
	return nil
}
func (k *CodeDbRepository) CallSetObjectByIndex(objectName string, indexName string, indexValue interface{}, objectValue interface{}) error {
	return nil
}
func (k *CodeDbRepository) CallSetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}, objectValue interface{}) error {
	return nil
}
func (k *CodeDbRepository) CallGetObject(objectName string) (interface{}, error) {
	return nil, nil
}
func (k *CodeDbRepository) CallGetObjectByIndex(objectName string, indexName string, indexValue interface{}) (interface{}, error) {
	return nil, nil
}
func (k *CodeDbRepository) CallGetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}) (interface{}, error) {
	return nil, nil
}
