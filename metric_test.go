package mackerelnullbridge_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	mackerelnullbridge "github.com/mashiike/mackerel-null-bridge"
	"github.com/stretchr/testify/require"
)

func TestFillConst(t *testing.T) {
	cases := []struct {
		from, to  int64
		values    []mackerel.MetricValue
		name      string
		fillValue interface{}
		expected  []*mackerel.MetricValue
	}{
		{
			from:      timeToEporch("2006-01-02T15:03:05Z"),
			to:        timeToEporch("2006-01-02T15:06:10Z"),
			name:      "hoge.fuga.piyo",
			fillValue: 0.0,
			expected: []*mackerel.MetricValue{
				{
					Name:  "hoge.fuga.piyo",
					Time:  timeToEporch("2006-01-02T15:04:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.piyo",
					Time:  timeToEporch("2006-01-02T15:05:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.piyo",
					Time:  timeToEporch("2006-01-02T15:06:00Z"),
					Value: 0.0,
				},
			},
		},
		{
			from:      timeToEporch("2006-01-02T15:04:05Z"),
			to:        timeToEporch("2006-01-02T15:08:10Z"),
			name:      "hoge.fuga.tora",
			fillValue: 0.0,
			values: []mackerel.MetricValue{
				{
					Name:  "hoge.fuga.tora",
					Time:  timeToEporch("2006-01-02T15:05:00Z"),
					Value: 60.0,
				},
				{
					Name:  "hoge.fuga.tora",
					Time:  timeToEporch("2006-01-02T15:07:00Z"),
					Value: 30.0,
				},
			},
			expected: []*mackerel.MetricValue{
				{
					Name:  "hoge.fuga.tora",
					Time:  timeToEporch("2006-01-02T15:06:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.tora",
					Time:  timeToEporch("2006-01-02T15:08:00Z"),
					Value: 0.0,
				},
			},
		},
		{
			from:      timeToEporch("2006-01-02T15:03:05Z"),
			to:        timeToEporch("2006-01-02T15:08:10Z"),
			name:      "hoge.fuga.tara",
			fillValue: 0.0,
			values: []mackerel.MetricValue{
				{
					Name:  "hoge.fuga.tara",
					Time:  timeToEporch("2006-01-02T15:05:00Z"),
					Value: 60.0,
				},
				{
					Name:  "hoge.fuga.tara",
					Time:  timeToEporch("2006-01-02T15:07:00Z"),
					Value: 30.0,
				},
			},
			expected: []*mackerel.MetricValue{
				{
					Name:  "hoge.fuga.tara",
					Time:  timeToEporch("2006-01-02T15:04:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.tara",
					Time:  timeToEporch("2006-01-02T15:06:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.tara",
					Time:  timeToEporch("2006-01-02T15:08:00Z"),
					Value: 0.0,
				},
			},
		},
		{
			from:      timeToEporch("2006-01-02T15:04:00Z"),
			to:        timeToEporch("2006-01-02T15:10:00Z"),
			name:      "hoge.fuga.kuma",
			fillValue: 0.0,
			values: []mackerel.MetricValue{
				{
					Name:  "hoge.fuga.kuma",
					Time:  timeToEporch("2006-01-02T15:05:00Z"),
					Value: 60.0,
				},
				{
					Name:  "hoge.fuga.kuma",
					Time:  timeToEporch("2006-01-02T15:08:00Z"),
					Value: 30.0,
				},
				{
					Name:  "hoge.fuga.kuma",
					Time:  timeToEporch("2006-01-02T15:10:00Z"),
					Value: 30.0,
				},
			},
			expected: []*mackerel.MetricValue{
				{
					Name:  "hoge.fuga.kuma",
					Time:  timeToEporch("2006-01-02T15:06:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.kuma",
					Time:  timeToEporch("2006-01-02T15:07:00Z"),
					Value: 0.0,
				},
				{
					Name:  "hoge.fuga.kuma",
					Time:  timeToEporch("2006-01-02T15:09:00Z"),
					Value: 0.0,
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			actual, err := mackerelnullbridge.FillConst(c.from, c.to, c.values, c.name, c.fillValue)
			require.NoError(t, err)
			require.EqualValues(t, c.expected, actual)
		})
	}
}

func timeToEporch(timeStr string) int64 {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}
	return t.Unix()
}
