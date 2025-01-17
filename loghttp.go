// +build go1.7

// Package loghttp provides automatic logging functionalities to http.Client.
package loghttp

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Transport implements http.RoundTripper. When set as Transport of http.Client, it executes HTTP requests with logging.
// No field is mandatory.
type Transport struct {
	Transport   http.RoundTripper
	LogRequest  func(req *http.Request)
	LogResponse func(resp *http.Response)
}

// THe default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

// Used if transport.LogRequest is not set.
var DefaultLogRequest = func(req *http.Request) {
	log.Printf("--> %s %s", req.Method, req.URL)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Run ioutil.ReadAll(resp.Body) failed")
	}
	log.Printf("<-- Header:\n %s\n", req.Header)
	log.Printf("<-- Body:\n %s\n", b)
}

// Used if transport.LogResponse is not set.
var DefaultLogResponse = func(resp *http.Response) {
	ctx := resp.Request.Context()
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Run ioutil.ReadAll(resp.Body) failed")
	}

	if _, ok := ctx.Value(ContextKeyRequestStart).(time.Time); ok {
		log.Printf("<-- %d %s", resp.StatusCode, resp.Request.URL)
		log.Printf("<-- Body:\n %s\n", respBody)
	} else {
		log.Printf("<-- %d %s", resp.StatusCode, resp.Request.URL)
	}
}

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), ContextKeyRequestStart, time.Now())
	req = req.WithContext(ctx)

	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *Transport) logRequest(req *http.Request) {
	if t.LogRequest != nil {
		t.LogRequest(req)
	} else {
		DefaultLogRequest(req)
	}
}

func (t *Transport) logResponse(resp *http.Response) {
	if t.LogResponse != nil {
		t.LogResponse(resp)
	} else {
		DefaultLogResponse(resp)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
