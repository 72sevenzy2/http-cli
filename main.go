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

func (h *HeaderFlags) String() string {
	return fmt.Sprint(*h)
}

func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func validate(args []string, bound int) error {
	if len(args) < bound {
		UsageMsg := errors.New("usage > main.go <URL> [-H key:value]")
		return fmt.Errorf("%s", UsageMsg.Error())
	}
	return nil
}

// func for adding headers
func addHeaders(req *http.Request, args HeaderFlags) error {
	for _, h := range args {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid input type %s", h)
		}

		// appending errors
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return nil
}

func main() {
	var headers HeaderFlags

	flag.Var(&headers, "H", "Header (key:value)")
	stream := flag.Bool("stream", false, "live response") // for streaming live response

	flag.Parse()

	if err := validate(flag.Args(), 2); err != nil {
		fmt.Println(errors.New(err.Error()))
		return
	}

	url := flag.Args()[0]

	req, err := http.NewRequest("GET", url, nil) // initialising an http request
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// add heders
	if err := addHeaders(req, headers); err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // timeout so that request doesnt last forever
	} // initialising the client

	// latency tracking
	start := time.Now()
	resp, err := client.Do(req)
	end := time.Since(start)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer resp.Body.Close() // important as not closing resp.Body would lead to performance issues + leaks, aswell as its apart of the ReadCloser interface so it has be closed.

	if *stream {
		var bf bytes.Buffer

		tee := io.TeeReader(resp.Body, &bf)

		_, err := io.Copy(os.Stdout, tee)
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
