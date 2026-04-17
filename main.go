package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	headers := flag.String("H", "", "Headers (key:value,key:value)")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("usage > main.go <URL> [-H key:value]")
		return
	}

	url := flag.Args()[0]

	req, err := http.NewRequest("GET", url, nil) // initialising an http request
	if err != nil {
		log.Fatal(err.Error())
	}

	if *headers != "" {
		headerPairs := strings.SplitSeq(*headers, ",") // split each header by a commar (SplitSeq then returns a string array containing the sorted headers)
		for h := range headerPairs {
			parts := strings.SplitN(h, ":", 2) // splitting each header by key:value pairs 
			if len(parts) == 2 {
				req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])) // trim the input if any redundant " " (spaces) are included when user sets headers
			}
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

	fmt.Println("\nbody:")

	var format bytes.Buffer // pretty printed body will be stored here before outputted

	err = json.Indent(&format, body, "", "  ")
	if err == nil {
		fmt.Println(format.String())
	} else {
		fmt.Println(string(body))
	}
}
