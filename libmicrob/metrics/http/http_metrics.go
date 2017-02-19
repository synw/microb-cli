package http_metrics

import (
	//"fmt"
	"net/http"
	"time"
	"github.com/tcnksm/go-httpstat"
	"github.com/synw/microb/libmicrob/datatypes"
)


func New(url string, server *datatypes.Server, server_processing int, total_time int) *datatypes.HttpRequestMetric {
	m := &datatypes.HttpRequestMetric{server, url, server_processing, total_time}
	return m
}

func GetRequestMetric(path string, server *datatypes.Server) (*datatypes.HttpRequestMetric, error) {
	var port string
	var metric *datatypes.HttpRequestMetric
	if server.Port != "" {
		port = ":"+server.Port
	}
	url := "http://"+server.Host+port+path
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
	server_processing := int(result.ServerProcessing/time.Millisecond)
	total_time := int(result.Total(t)/time.Millisecond)
	metric = New(path, server, server_processing, total_time)
	
	res.Body.Close()
	result.End(time.Now())

	return metric, nil
}
