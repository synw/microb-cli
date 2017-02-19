package http_metrics

import (
	//"fmt"
	"net/http"
	"time"
	"errors"
	"strings"
	"github.com/tcnksm/go-httpstat"
	"github.com/synw/microb/libmicrob/datatypes"
)


func New(url string, server *datatypes.Server, server_processing int, transport_time int, total_time int, status_code int) *datatypes.HttpRequestMetric {
	m := &datatypes.HttpRequestMetric{server, url, server_processing, transport_time, total_time, status_code}
	return m
}

func GetRequestMetric(path string, server *datatypes.Server) (*datatypes.HttpRequestMetric, error) {
	var port string
	var metric *datatypes.HttpRequestMetric
	if server.Port != "" {
		port = ":"+server.Port
	}
	var url string
	if strings.HasPrefix(path, "http://") {
		url = path
	} else if strings.HasPrefix(path, "/") == false {
		err := errors.New("Please start your url with / or http://")
		return metric, err
	} else {
		url = "http://"+server.Host+port+path
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
