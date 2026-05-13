package main

type key interface {
	string | int
}

type Data_storage[k key] map[k]string // using generics to support only types string and int as key name.