package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Log(v *http.Client, req *http.Request, bodyAllowed *bool) (time.Duration, *http.Response, error) {
	start := time.Now()
	resp, err := v.Do(req)
	if err != nil {
		return 0, nil, err
	}
	end := time.Since(start)

	// fmt.Println("visited to:", req.URL.Path)
	// fmt.Println("method:", req.Method)
	fmt.Printf("visited to %s, with method %s", req.URL.Path, req.Method)
	fmt.Println("full url:", req.URL)
	// request query
	if req.URL.RawQuery != "" {
		fmt.Println("with query:", req.URL.RawQuery)
	}

	fmt.Println("user details:")
	fmt.Println("client:", req.RemoteAddr)

	// exlude sensitive headers
	fmt.Println("request header details:")
	newHeaders := req.Header.Clone()
	newHeaders.Del("Authorization")
	fmt.Println(newHeaders) // then display

	// request body printing

	if *bodyAllowed {
		max := 1024 // 1 kb default (will add customisable max sizes later on)
		bodybytes, err := io.ReadAll(req.Body)
		if err != nil {
			if len(bodybytes) == 0 {
				fmt.Println("request does not contain any body.")
			} else {
				fmt.Println(err.Error())
			}
		}

		req.Body = io.NopCloser(bytes.NewBuffer(bodybytes))
		if len(bodybytes) > max {
			bodybytes = bodybytes[:max]
		} else {
			bodyprev := string(bodybytes)
			fmt.Println("request body:", bodyprev)
		}
	}

	return end, resp, nil
}
