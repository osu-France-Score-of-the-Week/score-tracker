package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type task struct {
	request  Request
	response chan taskResult
}

type taskResult struct {
	response Response
	err      error
}

type RequestQueue struct {
	tasks  chan task
	client *http.Client
}

func NewQueue(interval time.Duration) *RequestQueue {
	q := &RequestQueue{
		tasks:  make(chan task, 500),
		client: &http.Client{Timeout: 10 * time.Second},
	}

	go q.worker(interval)

	return q
}

func (q *RequestQueue) worker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for t := range q.tasks {
		resp, err := q.send(t.request)
		t.response <- taskResult{resp, err}
		<-ticker.C
	}
}

func (q *RequestQueue) send(req Request) (Response, error) {
	var body io.Reader
	if req.Body != nil {
		body = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return Response{}, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := q.client.Do(httpReq)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// Execute queues req and blocks until the queue has sent it, respecting the
// interval configured on q.
func (q *RequestQueue) Execute(req Request) (Response, error) {
	ch := make(chan taskResult, 1)
	q.tasks <- task{request: req, response: ch}
	res := <-ch
	return res.response, res.err
}

// ExecuteJSON queues req through q and decodes a successful JSON body into T,
// so callers get a typed result without writing status-check/decode boilerplate:
// `players, err := queue.ExecuteJSON[osuModels.UsersResponse](s.queue, req)`.
func ExecuteJSON[T any](q *RequestQueue, req Request) (T, error) {
	var out T

	resp, err := q.Execute(req)
	if err != nil {
		return out, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out, fmt.Errorf("request to %s failed with status %d: %s", req.URL, resp.StatusCode, string(resp.Body))
	}

	if err := json.Unmarshal(resp.Body, &out); err != nil {
		return out, err
	}

	return out, nil
}
