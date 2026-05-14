package main

import (
	"fmt"
)

func TestGet() {
	store := NewStore()
	_, ok := store.Set(1, "test")
	if ok {
		fmt.Println("data set successfully")
	} else {
		fmt.Println("could not save data")
	}

	val, exists := store.Get(1)
	if exists {
		fmt.Println("data exists:", val)
	} else {
		fmt.Println("couldnt retrieve data")
	}
}