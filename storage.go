// to hold user inputs during session mode - storing variables in maps to then be retrieved later using dynamic lookups.

package main

type key interface {
	string | int // only supports these 2 types for now
}

type data[k key] struct {
	data_storage map[k]string // using generics to support only types string and int as key name.4
}

// new db
func NewStore[k key]() *data[k] {
	return &data[k]{
		make(map[k]string),
	}
}

// utility get/set functions for data map:

// for strings
func (d *data[k]) Get(keyname k) (string, bool) {
	val, ok := d.data_storage[keyname]
	return val, ok
}

func (d *data[k]) Set(keyname k, value string) bool {
	// validate if value exists first
	if value != "" {
		d.data_storage[keyname] = value
		return true
	} else {
		return false // and then handle error from the module calling this method.
	}
}