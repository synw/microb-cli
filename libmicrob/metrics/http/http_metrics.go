package http_metrics

import (
	"fmt"
	"net/http"
	"time"
	"errors"
	"strings"
	"io/ioutil"
	"encoding/binary"
	"github.com/tcnksm/go-httpstat"
	"github.com/synw/microb/libmicrob/datatypes"
)


func New(url string, server *datatypes.Server, server_processing int, transport_time int, total_time int, status_code int) *datatypes.HttpRequestMetric {
	m := &datatypes.HttpRequestMetric{server, url, server_processing, transport_time, total_time, status_code}
	return m
}

func NewHttpResponse(server *datatypes.Server, url string, content string, size int, status_code int, duration *time.Duration) *datatypes.HttpResponse {
	r := &datatypes.HttpResponse{server, url, content, size, status_code, duration}
	return r	
}

func Get(path string, server *datatypes.Server ) (*datatypes.HttpResponse, error) {
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
    } else {
    	defer resp.Body.Close()
		// reads html as a slice of bytes
		bcontent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
	    	return r, err
	    }
		content = string(bcontent[:])
		size = binary.Size(bcontent)
	    }
	duration := time.Since(t)
	status_code := resp.StatusCode
	r = NewHttpResponse(server, url, content, size, status_code, &duration)
	fmt.Println(r)
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
