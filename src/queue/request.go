package queue

import "net/http"

type Request struct {
	Method  string
	URL     string
	Body    []byte
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}
