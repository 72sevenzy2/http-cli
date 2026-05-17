package main

import (
	"fmt"
	"errors"
	"net/http"
	"strings"
)

func (h *HeaderFlags) String() string { // gets called internally by go's flag pkg, (type flag.Value expects a String() and Set() func)
	return fmt.Sprint(*h)
}

func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func Validate(args []string) error {
	if len(args) < 2 {
		UsageMsg := errors.New("usage > main.go <URL> [-H key:value]")
		return fmt.Errorf("%s", UsageMsg.Error())
	}
	return nil
}

// func for adding headers
func AddHeaders(req *http.Request, args HeaderFlags) error {
	for _, h := range args {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid input type %s", h)
		}

		// appending headers
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return nil
}