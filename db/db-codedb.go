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

package db

func (d *DBService) StartUndoSession() {
	d.codeRepo.CallStartUndoSession(true)
}

func (d *DBService) CreatObjectIndex(objectName string, indexName string, indexJson string) error {
	return d.codeRepo.CallCreatObjectIndex(objectName, indexName, indexJson)
}
func (d *DBService) SetObject(objectName string, key string, objectValue string) error {
	return d.codeRepo.CallSetObject(objectName, key, objectValue)
}
func (d *DBService) GetObject(objectName string, key string) (string, error) {
	return d.codeRepo.CallGetObject(objectName, key)
}
func (d *DBService) GetObjectByIndex(objectName string, indexName string, indexValue interface{}) (string, error) {
	return d.codeRepo.CallGetObjectByIndex(objectName, indexName, indexValue)
}
func (d *DBService) GetAllObjectKeys(objectName string) ([]string, error) {
	return d.codeRepo.CallGetAllObjectKeys(objectName)
}
func (d *DBService) GetAllObjects(objectName string) ([]string, error) {
	return d.codeRepo.CallGetAllObjects(objectName)
}
func (d *DBService) GetAllObjectsSortByIndex(objectName string, indexName string) ([]string, error) {
	return d.codeRepo.CallGetAllObjectsSortByIndex(objectName, indexName)
}
func (d *DBService) DeleteObject(objectName string, key string) (string, error) {
	return d.codeRepo.CallDeleteObject(objectName, key)
}
func (d *DBService) Commit() error {
	return d.codeRepo.CallCommit()
}
func (d *DBService) Rollback() error {
	return d.codeRepo.CallRollback()
}
func (d *DBService) Reset() {
	//TODO
	d.codeRepo.CallStartUndoSession(false)
}
