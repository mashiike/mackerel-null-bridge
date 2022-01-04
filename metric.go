package mackerelnullbridge

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

type ServiceMetric struct {
	ServiceName  string
	MetricName   string
	DefaultValue interface{}
	DelaySeconds int64
	client       MackerelClient
}

const (
	fetchInterval  = int64(6 * time.Hour / time.Second)
	metricInterval = 60
)

func (m *ServiceMetric) FillConst(ctx context.Context, from int64, to int64) error {
	from = from - m.DelaySeconds
	to = to - m.DelaySeconds
	log.Printf("[info] fill null with %v for %s `%s` in the period %d to %d\n", m.DefaultValue, m.ServiceName, m.MetricName, from, to)
	values, err := m.fetchMetricValues(ctx, from, to)
	if err != nil {
		return err
	}
	postValues, err := FillConst(from, to, values, m.MetricName, m.DefaultValue)
	if err != nil {
		return err
	}
	log.Printf("[info] there are %d missing values for %s `%s` in the period %d to %d\n", len(postValues), m.ServiceName, m.MetricName, from, to)
	return m.client.PostServiceMetricValues(m.ServiceName, postValues)
}

func (m *ServiceMetric) fetchMetricValues(ctx context.Context, from int64, to int64) ([]mackerel.MetricValue, error) {
	ret := make([]mackerel.MetricValue, 0, (to-from)/metricInterval+1)
	for current := from; current <= to; current += fetchInterval {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		cFrom := current
		cTo := cFrom + fetchInterval
		if cTo >= to {
			cTo = to
		}
		values, err := m.client.FetchServiceMetricValues(m.ServiceName, m.MetricName, cFrom, cTo)
		if err != nil {
			return nil, err
		}
		ret = append(ret, values...)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Time < ret[j].Time
	})
	return ret, nil
}

func FillConst(from, to int64, values []mackerel.MetricValue, name string, fillValue interface{}) ([]*mackerel.MetricValue, error) {
	log.Printf("[debug] fill `%s` in %d~%d\n", name, from, to)
	n := len(values)
	from = (from/metricInterval)*metricInterval + metricInterval
	if n == 0 {
		log.Printf("[debug] `%s` has no metrics in %d~%d", name, from, to)
		return generateMetricValues(from, to, name, fillValue), nil
	}
	filledValues := make([]*mackerel.MetricValue, 0)
	cursorFrom := from
	var cursorTo int64
	for i := 0; i < n; i++ {
		if values[i].Name != "" && values[i].Name != name {
			return nil, errors.New("mismatch metric name")
		}
		cursorTo = values[i].Time
		if cursorTo-cursorFrom >= metricInterval {
			filledValues = append(filledValues, generateMetricValues(cursorFrom, cursorTo-metricInterval, name, fillValue)...)
		}
		cursorFrom = cursorTo + metricInterval
	}
	if to-cursorTo >= metricInterval {
		filledValues = append(filledValues, generateMetricValues(cursorTo+metricInterval, to, name, fillValue)...)
	}
	return filledValues, nil
}

func generateMetricValues(from, to int64, name string, value interface{}) []*mackerel.MetricValue {
	values := make([]*mackerel.MetricValue, 0, (to-from)/metricInterval+1)
	for current := from; current <= to; current += metricInterval {
		values = append(values, &mackerel.MetricValue{
			Name:  name,
			Time:  current,
			Value: value,
		})
	}
	log.Printf("[debug] generate new metric values `%s` in %d~%d, %d data points", name, from, to, len(values))
	return values
}
