package executor

import (
	"encoding/base64"
	"net/url"
	"time"

	"github.com/levigross/grequests"
)

type RestExec struct {
	url     string
	timeout time.Duration
}

func NewRestExec(url string, timeout time.Duration) *RestExec {
	return &RestExec{url: url, timeout: timeout}
}

type externalExecutionResponse struct {
	Returncode uint32 `json:"returncode"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	Version    string `json:"version"`
}

type cacheEntry struct {
	ExecResult
	Duration time.Duration
}

var cacheExecResult = map[string]cacheEntry{}

func (e *RestExec) Exec(code []byte, arg string, env interface{}) (ExecResult, error) {
	if e, ok := cacheExecResult[arg]; ok {
		time.Sleep(e.Duration)
		return e.ExecResult, nil
	}
	start := time.Now()
	executable := base64.StdEncoding.EncodeToString(code)
	resp, err := grequests.Post(
		e.url,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			JSON: map[string]interface{}{
				"executable": executable,
				"calldata":   arg,
				"timeout":    e.timeout.Milliseconds(),
				"env":        env,
			},
			RequestTimeout: e.timeout,
		},
	)

	if err != nil {
		urlErr, ok := err.(*url.Error)
		if !ok || !urlErr.Timeout() {
			return ExecResult{}, err
		}
		// Return timeout code
		return ExecResult{Output: []byte{}, Code: 111}, nil
	}

	if resp.Ok != true {
		return ExecResult{}, ErrRestNotOk
	}

	r := externalExecutionResponse{}
	err = resp.JSON(&r)

	if err != nil {
		return ExecResult{}, err
	}

	if r.Returncode == 0 {
		cacheDur := time.Since(start)
		entry := ExecResult{Output: []byte(r.Stdout), Code: 0, Version: r.Version}
		cacheExecResult[arg] = cacheEntry{entry, cacheDur}
		return entry, nil
	} else {
		return ExecResult{Output: []byte(r.Stderr), Code: r.Returncode, Version: r.Version}, nil
	}
}
