/**
 * @author mlabarinas
 */

package godog

import (
	"github.com/sschepens/datadog-go/statsd"
	"os"
)

const (
	ENDPOINT            string = "datadog:8125"
	DEFAULT_BUFFER_SIZE int    = 10000
)

type AwsDogClient struct{}

var client *statsd.Client

func (a *AwsDogClient) RecordSimpleMetric(metricName string, value float64, tags ...string) {
	client.Count(metricName, int64(value), getTags(tags...), 1)
}

func (a *AwsDogClient) RecordCompoundMetric(metricName string, value float64, tags ...string) {
	client.Gauge(metricName, value, getTags(tags...), 1)
}

func (a *AwsDogClient) RecordFullMetric(metricName string, value float64, tags ...string) {
	client.TimeInMilliseconds(metricName, value, getTags(tags...), 1)
}

func (a *AwsDogClient) RecordSimpleTimeMetric(metricName string, fn action, tags ...string) (interface{}, error) {
	time, result, error := takeTime(fn)

	client.Count(metricName, time, getTags(tags...), 1)

	return result, error
}

func (a *AwsDogClient) RecordCompoundTimeMetric(metricName string, fn action, tags ...string) (interface{}, error) {
	time, result, error := takeTime(fn)

	client.Gauge(metricName, float64(time), getTags(tags...), 1)

	return result, error
}

func (a *AwsDogClient) RecordFullTimeMetric(metricName string, fn action, tags ...string) (interface{}, error) {
	time, result, error := takeTime(fn)

	client.TimeInMilliseconds(metricName, float64(time), getTags(tags...), 1)

	return result, error
}

func getTags(tags ...string) []string {
	defaultTags := new(Tags).Add("platform", os.Getenv("PLATFORM")).Add("application", os.Getenv("APPLICATION")).Add("datacenter", os.Getenv("DATACENTER")).ToArray()

	defaultTags = append(defaultTags, tags...)

	return defaultTags
}

func init() {
	if os.Getenv("GO_ENVIRONMENT") == "production" {
		c, error := statsd.NewBuffered(ENDPOINT, DEFAULT_BUFFER_SIZE)

		if error != nil {
			panic(error)
		}

		client = c
		return
	}
}
