package executor

import (
	"encoding/base64"
	"net/url"
	"sync"
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

var cacheExecResult = sync.Map{} // map[string]cacheEntry{}

const yodaInitFile = "/home/panu/chain/yoda-init.txt"

func init() {
	// file, err := os.Open(yodaInitFile)
	// if err != nil {
	// 	panic(err)
	// }

	// scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanLines)

	// for scanner.Scan() {
	// 	text := scanner.Text()
	// 	msgs := strings.Split(text, "=")
	// 	k := msgs[0]
	// 	v := msgs[1]

	// 	d, err := hex.DecodeString(v)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	var e cacheEntry
	// 	if err := json.Unmarshal(d, &e); err != nil {
	// 		panic(err)
	// 	}

	// 	cacheExecResult.Store(k, e)
	// }
}

func (e *RestExec) Exec(code []byte, arg string, env interface{}) (ExecResult, error) {
	if v, ok := cacheExecResult.Load(arg); ok {
		e := v.(cacheEntry)
		time.Sleep(e.Duration)
		return e.ExecResult, nil
	}
	// fmt.Println(arg)
	// panic("not found why!")

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
		cacheExecResult.Store(arg, cacheEntry{entry, cacheDur})

		// f, err := os.OpenFile(yodaInitFile,
		// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	log.Println(err)
		// }
		// defer f.Close()
		// rawEntry, _ := json.Marshal(entry)
		// hexRawEntry := hex.EncodeToString(rawEntry)
		// if _, err := f.WriteString(fmt.Sprintf("%s=%s\n", arg, hexRawEntry)); err != nil {
		// 	log.Println(err)
		// }

		return entry, nil
	} else {
		return ExecResult{Output: []byte(r.Stderr), Code: r.Returncode, Version: r.Version}, nil
	}
}
