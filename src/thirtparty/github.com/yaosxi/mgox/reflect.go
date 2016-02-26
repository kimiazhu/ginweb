// Copyright 2015 The mgox Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgox

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

func getCollectionName(doc interface{}) string {
	if collectionName, ok := doc.(string); ok {
		return collectionName
	}
	var collectionName string
	mType := reflect.TypeOf(doc)
	mValue := reflect.ValueOf(doc)
	mType, mValue = getElem(mType, mValue)
	if mType.Kind() == reflect.Struct && mValue.IsValid() {
		field := mValue.FieldByName("CollectionName")
		if field.IsValid() && field.String() != "" {
			collectionName = field.String()
		}
	}
	if collectionName == "" && mType.Name() != "" {
		collectionName = strings.ToLower(mType.Name())
	}
	return collectionName
}

func getObjectId(doc interface{}) bson.ObjectId {
	mType := reflect.TypeOf(doc)
	mValue := reflect.ValueOf(doc)
	mType, mValue = getElem(mType, mValue)
	if mType.Kind() == reflect.Struct && mValue.IsValid() {
		field := mValue.FieldByName("Id")
		if field.IsValid() && field.Interface() != nil {
			if v, ok := field.Interface().(bson.ObjectId); ok {
				return v
			}
		}
	}
	return ""
}

func getElem(t reflect.Type, value reflect.Value) (reflect.Type, reflect.Value) {
	switch t.Kind() {
	case reflect.Ptr:
		//		fmt.Printf("\nDetected Ptr")
		return getElem(t.Elem(), value.Elem())
	case reflect.Slice:
		//		fmt.Printf("\nDetected Slice %d", value.Len())
		var val reflect.Value
		if value.Len() > 0 {
			val = value.Index(0)
		}
		return getElem(t.Elem(), val)
	case reflect.Array:
		//		fmt.Printf("\nDetected Array %d", value.Len())
		var val reflect.Value
		if value.Len() > 0 {
			val = value.Index(0)
		}
		return t.Elem(), val
	case reflect.Interface: // ????
		//		fmt.Printf("\nDetected Interface")
		return reflect.TypeOf(value), reflect.ValueOf(value)
	}
	return t, value
}

func IsSlice(v interface{}) bool {
	resultv := reflect.ValueOf(v)
	return resultv.Kind() == reflect.Ptr && resultv.Elem().Kind() == reflect.Slice
}

func GetValueLen(v interface{}) int {
	if IsSlice(v) {
		return reflect.ValueOf(v).Elem().Len()
	} else {
		return 1
	}
}
