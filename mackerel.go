package mackerelnullbridge

import (
	"log"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

type MackerelClient interface {
	FetchServiceMetricValues(serviceName string, metricName string, from int64, to int64) ([]mackerel.MetricValue, error)
	PostServiceMetricValues(serviceName string, metricValues []*mackerel.MetricValue) error
}

type DryRunMackerelClient struct {
	MackerelClient
}

func (c DryRunMackerelClient) PostServiceMetricValues(serviceName string, metricValues []*mackerel.MetricValue) error {
	log.Println("[notice] post service metrics values **dry run**")
	for _, v := range metricValues {
		log.Printf(`[notice] service=%s, data_point_time=%s, value={"name":%s, "time":%d, "value":%v}`, serviceName, time.Unix(v.Time, 0).Local(), v.Name, v.Time, v.Value)
	}
	return nil
}

func NewMackerelClient(apikey string, deploy bool) MackerelClient {
	client := mackerel.NewClient(apikey)
	if deploy {
		return client
	}
	return DryRunMackerelClient{
		MackerelClient: client,
	}
}
