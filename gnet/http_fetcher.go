// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package gnet

import (
	"fmt"
	"ghorgs/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Url     string
	Method  string
	Headers map[string]string
	Query   string
	Timeout time.Duration // in sec
}

func (r *Request) Execute() []byte {
	reqq := ""
	if r.Method == postMethod {
		reqq = r.Query
	}
	reqbody := strings.NewReader(reqq)
	req, err := http.NewRequest(r.Method, r.Url, reqbody)
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
		if utils.Debug.Verbose {
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
