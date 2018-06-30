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
 * file description:  none lock queue
 * @Author: eripi
 * @Date:   2018-1-05
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import "container/list"

type Queue struct {
	l *list.List
}

func NewQueue() *Queue {
	return &Queue{l: list.New()}
}

func (q *Queue) Pop() interface{} {
	if q.l.Front() != nil {
		return q.l.Remove(q.l.Front())
	}

	return nil
}

func (q *Queue) Push(data interface{}) {
	q.l.PushBack(data)
}

func (q *Queue) Length() int {
	return q.l.Len()
}
