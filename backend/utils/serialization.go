package utils

import (
	"encoding/json"
	"reflect"
)

func ToJSON(s interface{}) ([]byte, error) {
	return json.Marshal(s)
}

func FromJSON(data []byte, s interface{}) error {
    return json.Unmarshal(data, s)
}

func GetTypeName(v interface{}) string {
    t := reflect.TypeOf(v)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    return t.Name()
}
