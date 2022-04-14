package reflectutil

import (
	"encoding/json"
	"fmt"
)

type (
	// M is a map prepresent a struct information.
	M map[string]interface{}
)

// Parse convert a struct to a json object/map.
func Parse(v interface{}) (M, error) {
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
func (m M) Keys() []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// Values return values of the json object.
func (m M) Values() []interface{} {
	values := make([]interface{}, 0)
	for _, val := range m {
		values = append(values, val)
	}
	return values
}

// StringValues return values of the json object as strings.
func (m M) StringValues() []string {
	values := make([]string, 0)
	for _, val := range m {
		values = append(values, fmt.Sprintf("%v", val))
	}
	return values
}

// Sets set value of the keys to the given value.
func (m M) Sets(keys []string, value interface{}) {
	for _, key := range keys {
		m[key] = value
	}
}
