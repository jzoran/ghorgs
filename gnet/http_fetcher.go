//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package gnet

import (
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

// ResponseStatus holds the HTTP code and status resulting from an HTTP request.
type ResponseStatus struct {
	Code   int
	Status string
}

// Execute runs a given http request and returns resulting response body in bytes
// and  ResponseStatus (HTTP code and status).
func (r *Request) Execute() ([]byte, *ResponseStatus) {
	requestQuery := ""
	if r.Method == postMethod {
		requestQuery = r.Query
	}
	requestBody := strings.NewReader(requestQuery)
	req, err := http.NewRequest(r.Method, r.Url, requestBody)
	if err != nil {
		panic(err)
	}

	for key, header := range r.Headers {
		req.Header.Set(key, header)
	}

	var netClient = &http.Client{
		Timeout: time.Second * r.Timeout,
	}

	response, err := netClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	responseStatus := &ResponseStatus{response.StatusCode, response.Status}
	if responseStatus.Code != http.StatusOK {
		if utils.Debug.Verbose {
			log.Print(r.Query)
		}
		return nil, responseStatus
	}

	bbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return bbody, responseStatus
}
