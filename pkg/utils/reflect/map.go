package reflectutil

import (
	"encoding/json"
	"fmt"
)

type (
	// JsonObject is a map prepresent a struct information.
	JsonObject map[string]interface{}
)

// ToJsonObject convert a struct to a json object/map.
func ToJsonObject(v interface{}) (JsonObject, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Keys return keys of the json object.
func (jo JsonObject) Keys() []string {
	v := make([]string, 0)
	for k := range jo {
		v = append(v, k)
	}
	return v
}

// Values return values of the json object.
func (jo JsonObject) Values() []interface{} {
	v := make([]interface{}, 0)
	for _, vv := range jo {
		v = append(v, vv)
	}
	return v
}

// StringValues return values of the json object as strings.
func (jo JsonObject) StringValues() []string {
	v := make([]string, 0)
	for _, vv := range jo {
		v = append(v, fmt.Sprintf("%v", vv))
	}
	return v
}

// Sets set value of the keys to the given value.
func (jo JsonObject) Sets(keys []string, value interface{}) {
	for _, k := range keys {
		jo[k] = value
	}
}
