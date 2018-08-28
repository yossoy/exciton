package event

import (
	"encoding/json"
	"errors"
	"reflect"
)

// Value is hold event value like argument, return value, etc.
type Value interface {
	Encode() ([]byte, error)
	Decode(data interface{}) error
}

type jsonEncodedValue []byte

func (ev jsonEncodedValue) Encode() ([]byte, error) {
	return ev, nil
}
func (ev jsonEncodedValue) Decode(data interface{}) error {
	return json.Unmarshal(ev, data)
}

// NewJSONEncodedValue create json-encoded Value by interface{}.
func NewJSONEncodedValue(data interface{}) (Value, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonEncodedValue(b), nil
}

// NewJSONEncodedValueByEncodedBytes create json-encoded Value by json-encoded bytes.
func NewJSONEncodedValueByEncodedBytes(bytes []byte) Value {
	return jsonEncodedValue(bytes)
}

type valueEventParameter struct {
	value interface{}
}

func (ep valueEventParameter) Encode() ([]byte, error) {
	return json.Marshal(ep.value)
}
func (ep valueEventParameter) Decode(data interface{}) error {
	pdv := reflect.ValueOf(data)
	if pdv.IsNil() || pdv.Kind() != reflect.Ptr {
		return errors.New("invalid parameter type")
	}
	dv := pdv.Elem()
	if ep.value == nil {
		switch dv.Kind() {
		case reflect.Interface:
		case reflect.Ptr:
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		default:
			return errors.New("invaid type")
		}
	}
	vv := reflect.ValueOf(ep.value)
	if !vv.Type().AssignableTo(dv.Type()) {
		return errors.New("incompatible type")
	}
	dv.Set(vv)
	return nil
}

// NewValue create Value by interface{}
func NewValue(value interface{}) Value {
	return valueEventParameter{value: value}
}
