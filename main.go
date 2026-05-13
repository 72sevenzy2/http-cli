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
	"os"
	"strings"
	"time"
)

// custom type for composable header inputs
type HeaderFlags []string

func main() {
	var headers HeaderFlags

	flag.Var(&headers, "H", "Header (key:value)")
	stream := flag.Bool("stream", false, "live response") // for streaming live response
	method := flag.String("x", "GET", "http method")

	data := flag.String("d", "", "request data")

	flag.Parse()

	if err := Validate(flag.Args()); err != nil {
		fmt.Println(errors.New(err.Error()))
		return
	}

	url := flag.Args()[0]

	// validating whether method given is appropriate (only support for post and get for now)
	uc := strings.ToUpper(*method)
	var v string // holding actual method after validating if method is appropriate

	// simplified
	allowed := map[string]bool{
		"GET":  true,
		"POST": true,
	}

	if allowed[uc] {
		v = uc
	} else {
		fmt.Println("method not allowed")
		return
	}
	
	// body reader for post data
	var body io.Reader

	if *data != "" {
		body = strings.NewReader(*data)
	}

	req, err := http.NewRequest(v, url, body)

	// req, err := http.NewRequest("GET", url, nil) // initialising an http request
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// add heders
	if err := AddHeaders(req, headers); err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // timeout so that request doesnt last forever
	} // initialising the client

	// latency tracking
	// start := time.Now()
	// resp, err := client.Do(req)
	// end := time.Since(start)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// 	return
	// }

	// track latency
	end, resp, err := Log(client, req)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer resp.Body.Close() // important as not closing resp.Body would lead to performance issues + leaks, aswell as its apart of the ReadCloser interface so it has be closed.

	if *stream {
		_, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err.Error())
			return
		}

		// checking if request failed if yes then log
		if resp.StatusCode >= 400 { // anything over 400 means request wasnt successful
			fmt.Println("request failed with status:", resp.Status)
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
}
