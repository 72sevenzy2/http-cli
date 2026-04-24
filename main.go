package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// custom type for composable header inputs
type HeaderFlags []string

func (h *HeaderFlags) String() string {
	return fmt.Sprint(*h)
}

func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func validate(args []string, bound int) error {
	if len(args) < bound {
		UsageMsg := errors.New("usgae > main.go <URL> [-H key:value]")
		return fmt.Errorf("%s", UsageMsg.Error())
	}
	return nil
}

func main() {
	var headers HeaderFlags

	flag.Var(&headers, "H", "Header (key:value)")

	flag.Parse()

	if err := validate(flag.Args(), 2); err != nil {
		fmt.Println(err.Error())
		return
	}

	url := flag.Args()[0]

	req, err := http.NewRequest("GET", url, nil) // initialising an http request
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(headers) > 0 {
		for _, h := range headers { // headers returns a string array
			parts := strings.SplitN(h, ":", 2) // splitting each header by key:value pairs
			if len(parts) != 2 {
				log.Fatalf("invalid input type %s", h)
				return
			}
			req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])) // trimSpace removes any un-needed spaces in the input if present.
		}
	}

	client := &http.Client{} // initialising the client

	// latency tracking
	start := time.Now()
	resp, err := client.Do(req)
	end := time.Since(start)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close() // important as not closing resp.Body would lead to performance issues + leaks, aswell as its apart of the ReadCloser interface so it has be closed.

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	// outputting

	fmt.Println("status:", resp.Status)
	fmt.Println("latency:", end) // printing latency

	fmt.Println("\nheaders:")
	for k, v := range resp.Header {
		fmt.Println(k+":", v) // key:value output style for headers
	}

	fmt.Println("\nresponse size:")
	fmt.Println(len(body), "bytes")

	fmt.Println("\nbody:")

	var format bytes.Buffer // pretty printed body will be stored here before outputted

	err = json.Indent(&format, body, "", "  ")
	if err == nil {
		fmt.Println(format.String())
	} else {
		fmt.Println(string(body))
	}
}
