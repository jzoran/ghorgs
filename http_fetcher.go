package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Url     string
	Headers map[string]string
	Query   string
	Timeout time.Duration // in sec
}

func (r *Request) fetch() []byte {
	reqbody := strings.NewReader(r.Query)
	req, err := http.NewRequest("POST", r.Url, reqbody)
	if err != nil {
		panic(err)
	}

	for key, header := range r.Headers {
		req.Header.Set(key, header)
	}

	var netClient = &http.Client{
		Timeout: time.Second * r.Timeout,
	}

	resp, err := netClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if debug.Verbose {
			log.Print(r.Query)
		}
		panic(fmt.Sprintf("HttpResponse: %d", resp.StatusCode))
	}

	bbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return bbody
}
