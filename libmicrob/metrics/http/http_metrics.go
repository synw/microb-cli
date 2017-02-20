package http_metrics

import (
	"net/http"
	"time"
	"errors"
	"strings"
	"io/ioutil"
	"math/rand"
	"encoding/binary"
	"github.com/tcnksm/go-httpstat"
	"github.com/synw/microb/libmicrob/datatypes"
)


// constructors
func New(url string, server *datatypes.Server, server_processing int, transport_time int, total_time int, status_code int) *datatypes.HttpRequestMetric {
	m := &datatypes.HttpRequestMetric{server, url, server_processing, transport_time, total_time, status_code}
	return m
}

func NewHttpResponse(server *datatypes.Server, url string, content string, size int, status_code int, duration time.Duration) *datatypes.HttpResponse {
	r := &datatypes.HttpResponse{server, url, content, size, status_code, duration}
	return r	
}

func NewStressReport(server *datatypes.Server, num_requests int, size int, duration time.Duration) *datatypes.HttpStressReport {
	r := &datatypes.HttpStressReport{server, num_requests, size, duration}
	return r
}

// methods
func Stress(server *datatypes.Server, interval int, workers int, limit int) (*datatypes.HttpStressReport, error) {
	var duration time.Duration
	var size int
	
	d := time.Duration(interval)*time.Millisecond
	nw := 0
	i := 0
	for _ = range time.Tick(d) {
		nw = 0
		for nw < workers {
			url := "/"
			s := rand.Intn(500)
			time.Sleep(time.Duration(s)*time.Millisecond)
			resp, _ := Get(url, server)
		    size = size+resp.Size
		    duration = duration+resp.Duration
			nw = nw+1
			i++
		}
		if i >= limit {
			break
		}
	}
	report := NewStressReport(server, limit, size, duration)
	return report, nil
}

func Get(path string, server *datatypes.Server) (*datatypes.HttpResponse, error) {
	var r *datatypes.HttpResponse
	t := time.Now()
	var content string
	var size int
	url, err := getUrlFromPath(path, server)
	if err != nil {
		return r, err
	}
	resp, err := http.Get(url)
    if err != nil {
    	return r, err
    }
	defer resp.Body.Close()
	bcontent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    	return r, err
    }
	content = string(bcontent[:])
	size = binary.Size(bcontent)
	duration := time.Since(t)
	status_code := resp.StatusCode
	r = NewHttpResponse(server, url, content, size, status_code, duration)
    return r, nil
}

func GetRequestMetric(path string, server *datatypes.Server) (*datatypes.HttpRequestMetric, error) {
	var metric *datatypes.HttpRequestMetric
	url, err := getUrlFromPath(path, server)
	if err != nil {
		return metric, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return metric, err
	}
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)
	
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return metric, err
	}
	
	t := time.Now()
	sc := res.StatusCode
	server_processing := int(result.ServerProcessing/time.Millisecond)
	total_time := int(result.Total(t)/time.Millisecond)
	transport_time := total_time-server_processing
	metric = New(path, server, server_processing, transport_time, total_time, sc)
	
	res.Body.Close()
	result.End(time.Now())

	return metric, nil
}

func getUrlFromPath(path string, server *datatypes.Server) (string, error) {
	var url string
	if strings.HasPrefix(path, "http://") {
		url = path
		return url, nil
	} else if strings.HasPrefix(path, "/") == false {
		err := errors.New("Please start your url with / or http://")
		return url, err
	} else {
		url = "http://"+server.Host+":"+server.Port+path
	}
	return url, nil
}
