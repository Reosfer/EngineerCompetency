package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Options :
type Options struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Body        []byte            `json:"body"`
	Headers     map[string]string `json:"headers"`
	Timeout     time.Duration
	ContentType string `json:"content_type"`
}

// Response :
type Response struct {
	StatusCode int    `json:"status-code"`
	Body       []byte `json:"body"`
	Error      error  `json:"error"`
}

func GET(opt *Options) Response {
	res := <-request(opt, "GET")
	return res
}

func POST(opt *Options) Response {
	res := <-request(opt, "POST")
	return res
}

func PUT(opt *Options) Response {
	res := <-request(opt, "PUT")
	return res
}

func DELETE(opt *Options) Response {
	res := <-request(opt, "DELETE")
	return res
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Request :
func request(opt *Options, method string) <-chan Response {
	res := make(chan Response)
	go func() {
		defer Recover("rest http")
		defer close(res)
		var rsp *http.Response
		var e error
		c := http.Client{
			Timeout:   opt.Timeout,
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		}
		logrus.Debugf("http request body : %s", opt.Body)
		clientReqHeader := http.Header{}
		for k, v := range opt.Headers {
			clientReqHeader.Add(k, v)
		}
		if opt.ContentType == "" {
			clientReqHeader.Add("Content-Type", "application/json")
		} else {
			clientReqHeader.Add("Content-Type", opt.ContentType)
		}
		reqObj := http.Request{}
		reqObj.Method = method
		reqObj.URL, _ = url.Parse(opt.URL)
		reqObj.Header = clientReqHeader
		reqObj.Body = ioutil.NopCloser(bytes.NewBuffer(opt.Body))

		rsp, e = c.Do(&reqObj)
		if e != nil {
			logrus.Debugf("error when creating http request %s", e.Error())
			res <- Response{Error: fmt.Errorf("failed to create new request")}
			return
		}
		defer rsp.Body.Close()
		body, e := ioutil.ReadAll(rsp.Body)
		res <- Response{StatusCode: rsp.StatusCode, Body: body}
	}()
	return res
}

func GetStatusCode(err error, statusCodeDefault int) int {
	var statusCode int
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(strings.ToLower(err.Error()), "already") {
		statusCode = http.StatusConflict
	} else if strings.Contains(strings.ToLower(err.Error()), "expired") {
		statusCode = http.StatusGone
	} else {
		statusCode = http.StatusInternalServerError
	}

	if statusCodeDefault > 0 {
		statusCode = statusCodeDefault
	}

	return statusCode
}

func RequestId() string {
	//make requestId from package uuid
	return uuid.New().String()
}
