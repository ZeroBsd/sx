// SPDX-License-Identifier: 0BSD
package sx

import "encoding/json"

// Tries to convert object to a compact/dense json string
func ToJson(object any) Result[string] {
	var bytes, err = json.Marshal(object)
	if err != nil {
		return NewResultFromError[string](err)
	}
	return NewResultFrom(string(bytes))
}

// Tries to convert object to a pretty-printed json string
//
// recommended indentationStrings are "\t" or multiple spaces
func ToJsonPretty(object any, indentationString string) Result[string] {
	var bytes, err = json.MarshalIndent(object, "", indentationString)
	if err != nil {
		return NewResultFromError[string](err)
	}
	return NewResultFrom(string(bytes))
}

// Tries to convert a json string into the specified object
func FromJson[T any](jsonString string) Result[T] {
	var object T
	var err = json.Unmarshal([]byte(jsonString), &object)
	if err != nil {
		var e = err.Error()
		_ = e
		return NewResultFromError[T](err)
	}
	return NewResultFrom(object)
}

func FromJsonInterface[T any](jsonString string, object *T) Result[T] {
	var err = json.Unmarshal([]byte(jsonString), object)
	if err != nil {
		var e = err.Error()
		_ = e
		return NewResultFromError[T](err)
	}
	return NewResultFrom(*object)
}

func (m hashMapImpl[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Map)
}

func (m hashMapImpl[K, V]) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &m.Map)
}
