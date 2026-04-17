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
		fmt.Println("usage > http-tester <URL> [-H key:value]")
		return
	}

	url := flag.Args()[0]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if *headers != "" {
		headerPairs := strings.Split(*headers, ",")
		for _, h := range headerPairs {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	client := &http.Client{}

	// latency tracking
	start := time.Now()
	resp, err := client.Do(req)
	end := time.Since(start)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// outputting

	fmt.Println("status:", resp.Status)
	fmt.Println("latency:", end)

	fmt.Println("\nheaders:")
	for k, v := range resp.Header {
		fmt.Println(k+":", v)
	}

	fmt.Println("\nbody:")

	var format bytes.Buffer

	err = json.Indent(&format, body, "", "  ")
	if err == nil {
		log.Fatal(format.String())
	} else {
		fmt.Println(string(body))
	}
}
