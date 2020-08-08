package mango

import (
	"reflect"
	"strings"

	"github.com/dsmontoya/utils/reflectutils"
	"github.com/dsmontoya/utils/strutils"
	"github.com/jinzhu/inflection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var clientKey = &mongo.Client{}
var keyConnection *Connection

func getCollection(i interface{}) string {
	t := reflect.TypeOf(i)
	ts := t.String()
	if ts == "string" {
		return i.(string)
	}
	split := strings.Split(ts, ".")
	name := split[len(split)-1]
	plural := inflection.Plural(name)
	b := []byte(plural)
	b[0] += 'a' - 'A'
	return string(b)
}

func arrayToBsonA(v reflect.Value) bson.A {
	l := v.Len()
	a := bson.A{}
	for i := 0; i < l; i++ {
		itemValue := v.Index(i)
		deepValue := reflectutils.DeepValue(itemValue)
		item := valueToBson(deepValue)
		a = append(a, item)
	}
	return a
}

func structToBsonDoc(v reflect.Value) bson.D {
	doc := bson.D{}
	n := v.NumField()
	t := v.Type()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		fieldName := field.Name
		if strings.Contains(string(field.Tag), `bson:"-"`) {
			continue
		}
		newFieldName := strutils.ToSnakeCase(fieldName)
		fieldValue := v.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}
		// TODO: this looks ugly
		if ft := fieldValue.Type(); ft.String() == "mango.Document" {
			// TODO: document fields
			continue
		}
		deepValue := reflectutils.DeepValue(fieldValue)
		newFieldValue := valueToBson(deepValue)
		element := bson.E{newFieldName, newFieldValue}
		doc = append(doc, element)
	}
	return doc
}

func toBsonDoc(model interface{}) bson.D {
	v := reflectutils.DeepValue(reflect.ValueOf(model))
	return valueToBson(v).(bson.D)
}

func valueToBson(v reflect.Value) interface{} {
	k := v.Kind()
	switch k {
	case reflect.Invalid:
		return bson.D{}
	case reflect.Struct:
		return structToBsonDoc(v)
	case reflect.Array, reflect.Slice:
		return arrayToBsonA(v)
	default:
		return v.Interface()
	}
	return bson.D{}
}
