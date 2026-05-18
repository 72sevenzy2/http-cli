// to hold user inputs during session mode - storing variables in maps to then be retrieved later using dynamic lookups.

package main

import (
	"errors"
	"strconv"
)

type Data struct {
	data_storage map[string]string // using generics to support only types string and int as key name.
}

// new db
func NewStore() *Data {
	return &Data{
		make(map[string]string),
	}
}

// normalize key types to string
func Normalize(keyname any) (string, error) {
	switch v := any(keyname).(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		errms := errors.New("invalid type: consider only string or int.")
		return "", errms
	}
}

// utility get/set functions for data map:

// for strings
func (d *Data) Get(keyname any) (string, bool, error) {
	newKey, err := Normalize(keyname)
	if err != nil {
		return "", false, err
	}

	val, ok := d.data_storage[newKey]
	return val, ok, nil
}

func (d *Data) Set(keyname any, value string) (error, bool) {
	if value == "" {
		errM := errors.New("please include a value")
		return errM, false
	}

	newK, err := Normalize(keyname)
	if err != nil {
		return err, false
	}
	d.data_storage[newK] = value
	return nil, true
}
