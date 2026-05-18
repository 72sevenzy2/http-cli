package main

import (
	"bufio"
	"fmt"
	"strings"
)

// start a interactive session
func StartSession(b *bufio.Scanner, store *Data) {
	fmt.Println("session started.")
	for {
		fmt.Print(">")
		b.Scan() // take input

		input := b.Text()
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue // skip current iteration if no input
		}

		upperInput := strings.ToUpper(parts[0]) // "VAR", "GET", "DEL", "EXIT"

		switch upperInput {
		case "VAR":
			err, ok := store.Set(parts[1], parts[2])
			if err != nil && !ok {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("successful")
			continue
		case "GET":
			val, ok, err := store.Get(parts[1])
			if err != nil && !ok {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println(val)
		case "EXIT":
			fmt.Println("exiting will remove saved variables.")
			return
		}

	}
}
