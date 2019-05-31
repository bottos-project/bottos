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
 * file description:  bpl decode
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */

package bpl

import (
	"fmt"
	"io"
	"math/big"
	"reflect"
)

type DecodeContext struct {
	r         io.Reader
	t         uint8
	ext       uint8
	size      uint16
	stopField string
	stoped    bool
	rootValue interface{}
}

//DecoderReader function type of the decoder
type DecoderReader func(reflect.Value, *DecodeContext) error

//Decoder interface for customization
type Decoder interface {
	DecodeBPL(io.Reader) error
}

var (
	decoderInterface = reflect.TypeOf(new(Decoder)).Elem()
)

//Decode decodes byte stream to struct, slice/array or other basic types
func Decode(r io.Reader, v interface{}, stopField string) error {
	if v == nil {
		return fmt.Errorf("bpl decode: nil pointer")
	}

	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("bpl decode: need a pointer")
	}
	if rv.IsNil() {
		return fmt.Errorf("bpl decode: nil pointer")
	}

	decoder, err := getDecoder(t.Elem())
	if err != nil {
		return err
	}

	ctx := newDecodeContext(r, stopField)
	if err = ctx.readHeader(); err != nil {
		return err
	}
	return decoder(rv.Elem(), ctx)
}

func getDecoder(t reflect.Type) (DecoderReader, error) {
	kind := t.Kind()
	switch {
	case t.AssignableTo(reflect.PtrTo(bigInt)):
		return decodeBigInt, nil
	case t.AssignableTo(bigInt):
		return decodeBigIntNoPtr, nil
	case kind == reflect.Bool:
		return decodeBool, nil
	case kind == reflect.Uint8:
		return decodeUint8, nil
	case kind == reflect.Uint16:
		return decodeUint16, nil
	case kind == reflect.Uint32:
		return decodeUint32, nil
	case kind == reflect.Uint64:
		return decodeUint64, nil
	case kind == reflect.String:
		return decodeString, nil
	case kind == reflect.Slice && t.Elem().Kind() == reflect.Uint8:
		return decodeBytes, nil
	case kind == reflect.Array && t.Elem().Kind() == reflect.Uint8:
		return decodeByteArray, nil
	case kind == reflect.Slice || kind == reflect.Array:
		return makeArrayDecoder(t)
	case kind == reflect.Struct:
		return makeStructDecoder(t)
	case kind == reflect.Ptr:
		return makePtrDecoder(t)
	default:
		return nil, fmt.Errorf("bpl decode: type %v, kind %v not support", t, kind)
	}
}

func newDecodeContext(r io.Reader, stopfield string) *DecodeContext {
	ctx := &DecodeContext{r: r, stopField: stopfield, stoped: false}
	return ctx
}

func (ctx *DecodeContext) reset() {
	ctx.t = 0
	ctx.size = 0
	ctx.ext = 0
}

func (ctx *DecodeContext) readHeader() error {
	ctx.reset()

	t, err := ReadByte(ctx.r)
	if err != nil {
		return err
	}

	ctx.t = t
	switch ctx.t {
	case NIL:
	case FALSE:
	case TRUE:
	case UINT8:
	case UINT16:
	case UINT32:
	case UINT64:
	case BIN16:
		ctx.size, _, err = ReadUint16(ctx.r)
		if err != nil {
			return err
		}
	case STR16:
		ctx.size, _, err = ReadUint16(ctx.r)
		if err != nil {
			return err
		}
	case ARRAY16:
		ctx.size, _, err = ReadUint16(ctx.r)
		if err != nil {
			return err
		}
	case EXT16:
		ctx.size, _, err = ReadUint16(ctx.r)
		if err != nil {
			return err
		}
		ctx.ext, err = ReadByte(ctx.r)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("bpl decode: unknown type identifier %X", ctx.t)
	}

	return nil
}

func (ctx *DecodeContext) readUint() (uint64, error) {
	var err error
	switch ctx.t {
	case UINT8:
		v, err := ReadByte(ctx.r)
		if err == nil {
			return uint64(v), nil
		}
	case UINT16:
		v, _, err := ReadUint16(ctx.r)
		if err == nil {
			return uint64(v), nil
		}
	case UINT32:
		v, _, err := ReadUint32(ctx.r)
		if err == nil {
			return uint64(v), nil
		}
	case UINT64:
		v, _, err := ReadUint64(ctx.r)
		if err == nil {
			return uint64(v), nil
		}
	}

	return uint64(0), err
}

func (ctx *DecodeContext) readRaw() ([]byte, error) {
	val := make([]byte, ctx.size)
	n, e := ctx.r.Read(val)
	if e == nil && n == int(ctx.size) {
		return val, nil
	}
	return []byte{}, nil
}

func decodeBigIntNoPtr(v reflect.Value, ctx *DecodeContext) error {
	return decodeBigInt(v.Addr(), ctx)
}

func decodeBigInt(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != EXT16 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, EXT16)
	}
	if ctx.ext != EXT_BIGINT {
		return fmt.Errorf("bpl decode: decode exttype %X, expected %X", ctx.ext, EXT_BIGINT)
	}
	val, err := ctx.readRaw()
	if err != nil {
		return fmt.Errorf("bpl decode: decode big.Int fail")
	}
	i := v.Interface().(*big.Int)
	if i == nil {
		i = new(big.Int)
		v.Set(reflect.ValueOf(i))
	}
	i.SetBytes(val)
	return nil
}

func decodeBool(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t == TRUE {
		v.SetBool(true)
	} else if ctx.t == FALSE {
		v.SetBool(false)
	} else {
		return fmt.Errorf("bpl decode: decode type %X, expected %X or %X", ctx.t, TRUE, FALSE)
	}

	return nil
}

func decodeUint8(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != UINT8 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, UINT8)
	}
	val, err := ctx.readUint()
	if err != nil {
		return fmt.Errorf("bpl decode: decode uint8 fail")
	}
	v.SetUint(val)
	return nil
}

func decodeUint16(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != UINT16 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, UINT16)
	}
	val, err := ctx.readUint()
	if err != nil {
		return fmt.Errorf("bpl decode: decode uint16 fail")
	}
	v.SetUint(val)
	return nil
}

func decodeUint32(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != UINT32 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, UINT32)
	}
	val, err := ctx.readUint()
	if err != nil {
		return fmt.Errorf("bpl decode: decode uint32 fail")
	}
	v.SetUint(val)
	return nil
}

func decodeUint64(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != UINT64 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, UINT64)
	}
	val, err := ctx.readUint()
	if err != nil {
		return fmt.Errorf("bpl decode: decode uint64 fail")
	}
	v.SetUint(val)
	return nil
}

func decodeString(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != STR16 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, STR16)
	}
	val, err := ctx.readRaw()
	if err != nil {
		return fmt.Errorf("bpl decode: decode string fail")
	}
	v.SetString(string(val))
	return nil
}

func decodeBytes(v reflect.Value, ctx *DecodeContext) error {
	if ctx.t != BIN16 {
		return fmt.Errorf("bpl decode: decode type %X, expected %X", ctx.t, BIN16)
	}
	val, err := ctx.readRaw()
	if err != nil {
		return fmt.Errorf("bpl decode: decode bin16 fail")
	}
	v.SetBytes(val)
	return nil
}

func decodeByteArray(v reflect.Value, ctx *DecodeContext) error {
	val, err := ctx.readRaw()
	if err != nil {
		return fmt.Errorf("bpl decode: decode bin16 fail")
	}

	vlen := v.Len()
	if len(val) != vlen {
		return fmt.Errorf("bpl decode: wrong array size")
	}

	slice := v.Slice(0, vlen).Interface().([]byte)
	copy(slice, val)

	return nil
}

func decodeArray(v reflect.Value, ctx *DecodeContext, elemDecoder DecoderReader) error {
	vlen := v.Len()
	if vlen != int(ctx.size) {
		return fmt.Errorf("bpl decode: wrong array size %v, expected %v", int(ctx.size), vlen)
	}
	i := 0
	for ; i < vlen; i++ {
		if err := elemDecoder(v.Index(i), ctx); err != nil {
			return fmt.Errorf("bpl decoder: array decode error")
		}
	}
	if i < vlen {
		return fmt.Errorf("bpl decode: array size array")
	}
	return nil
}

func decodeSlice(v reflect.Value, ctx *DecodeContext, elemDecoder DecoderReader) error {
	asize := int(ctx.size)
	if asize == 0 {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		return nil
	}

	i := 0
	for ; i < asize; i++ {
		if i >= v.Cap() {
			newcap := v.Cap() + v.Cap()/2
			if newcap < 4 {
				newcap = 4
			}
			newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
			reflect.Copy(newv, v)
			v.Set(newv)
		}
		if i >= v.Len() {
			v.SetLen(i + 1)
		}

		if err := ctx.readHeader(); err != nil {
			return err
		}

		if err := elemDecoder(v.Index(i), ctx); err != nil {
			return fmt.Errorf("bpl decode: slice decode error")
		}
	}
	return nil
}

func makeArrayDecoder(t reflect.Type) (DecoderReader, error) {
	etype := t.Elem()
	elemDecoder, err := getDecoder(etype)
	if err != nil {
		return nil, err
	}

	var dec DecoderReader
	switch {
	case t.Kind() == reflect.Array:
		dec = func(val reflect.Value, ctx *DecodeContext) error {
			return decodeArray(val, ctx, elemDecoder)
		}
	default: // t.Kind() == reflect.Slice
		dec = func(val reflect.Value, ctx *DecodeContext) error {
			return decodeSlice(val, ctx, elemDecoder)
		}
	}
	return dec, nil
}

func makeStructDecoder(t reflect.Type) (DecoderReader, error) {
	fields := []structField{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, structField{f.Type, i, false})
	}

	rule, hasRule := ignoreRuleMap[t.Name()]

	dec := func(val reflect.Value, ctx *DecodeContext) (err error) {
		if len(fields) < int(ctx.size) {
			return fmt.Errorf("bpl decode: struct feild num mismatch, num %v, expected %v", len(fields), int(ctx.size))
		}
		if hasRule && ctx.rootValue == nil {
			ctx.rootValue = val.Interface()
		}
		for _, f := range fields {
			if hasRule && rule(t.Field(f.index), f.index, val.Interface(), ctx.rootValue) { //ignore
				continue
			}
			decoder, err := getDecoder(f.t)
			if err != nil {
				return err
			}

			if err := ctx.readHeader(); err != nil {
				return err
			}
			err = decoder(val.Field(f.index), ctx)
			if err != nil {
				return err
			}
			if ctx.stoped {
				break
			}
			if len(ctx.stopField) > 0 && t.Field(f.index).Name == ctx.stopField {
				ctx.stoped = true
				break
			}
		}
		return nil
	}
	return dec, nil
}

func makePtrDecoder(t reflect.Type) (DecoderReader, error) {
	etype := t.Elem()
	decoder, err := getDecoder(etype)
	if err != nil {
		return nil, err
	}
	dec := func(val reflect.Value, ctx *DecodeContext) (err error) {
		newval := val
		if val.IsNil() {
			newval = reflect.New(etype)
		}

		if ctx.t == NIL {
			val.Set(reflect.Zero(t))
			return nil
		}

		if err = decoder(newval.Elem(), ctx); err == nil {
			val.Set(newval)
		}
		return err
	}
	return dec, nil
}
