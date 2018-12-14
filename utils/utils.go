package utils

import (
	"github.com/satori/go.uuid"
	"reflect"
)

func RandToken() string {
	return uuid.Must(uuid.NewV4()).String()
}

func Map(in interface{}, fn func(interface{}) interface{}) interface{} {
	val := reflect.ValueOf(in)
	out := make([]interface{}, val.Len())

	for i := 0; i < val.Len(); i++ {
		out[i] = fn(val.Index(i).Interface())
	}

	return out
}
