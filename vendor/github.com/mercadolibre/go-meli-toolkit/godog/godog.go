/**
 * @author mlabarinas
 */

package godog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Client interface {
	RecordSimpleMetric(metricName string, value float64, tags ...string)
	RecordCompoundMetric(metricName string, value float64, tags ...string)
	RecordFullMetric(metricName string, value float64, tags ...string)
	RecordSimpleTimeMetric(metricName string, fn action, tags ...string) (interface{}, error)
	RecordCompoundTimeMetric(metricName string, fn action, tags ...string) (interface{}, error)
	RecordFullTimeMetric(metricName string, fn action, tags ...string) (interface{}, error)
}

type action func() (interface{}, error)

type Tags struct {
	values map[string]string
}

func (t *Tags) Add(key string, value string) *Tags {
	t.init()

	if strings.TrimSpace(key) != "" && strings.TrimSpace(value) != "" {
		t.values[key] = value
	}

	return t
}

func (t *Tags) Remove(key string) *Tags {
	t.init()

	delete(t.values, key)

	return t
}

func (t *Tags) ToArray() []string {
	t.init()

	tags := make([]string, 0)

	for k, v := range t.values {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}

	return tags
}

func (t *Tags) init() {
	if t.values == nil {
		t.values = make(map[string]string)
	}
}

var instance Client

func RecordSimpleMetric(metricName string, value float64, tags ...string) {
	instance.RecordSimpleMetric(metricName, value, tags...)
}

func RecordCompoundMetric(metricName string, value float64, tags ...string) {
	instance.RecordCompoundMetric(metricName, value, tags...)
}

func RecordFullMetric(metricName string, value float64, tags ...string) {
	instance.RecordFullMetric(metricName, value, tags...)
}

func RecordSimpleTimeMetric(metricName string, fn action, tags ...string) (interface{}, error) {
	return instance.RecordSimpleTimeMetric(metricName, fn, tags...)
}

func RecordCompoundTimeMetric(metricName string, fn action, tags ...string) (interface{}, error) {
	return instance.RecordCompoundTimeMetric(metricName, fn, tags...)
}

func RecordFullTimeMetric(metricName string, fn action, tags ...string) (interface{}, error) {
	return instance.RecordFullTimeMetric(metricName, fn, tags...)
}

func takeTime(fn action) (int64, interface{}, error) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	result, error := fn()

	end := time.Now().UnixNano() / int64(time.Millisecond)

	return (end - start), result, error
}

func RecordApiCallMetric(target string,initTime time.Time,status string){
	benchTime := time.Now().Sub(initTime)

	RecordSimpleMetric("go.api_call.request", 1, new(Tags).Add("target_id", target).ToArray()...)
	RecordSimpleMetric("go.api_call.status", 1, new(Tags).Add("target_id", target).Add("status", status).ToArray()...)
	RecordCompoundMetric("go.api_call.time", benchTime.Seconds()*1000, new(Tags).Add("target_id", target).ToArray()...)
}

func init() {
	if os.Getenv("DATACENTER") == "AWS" {
		instance = new(AwsDogClient)

	} else {
		instance = new(OpenStackDogClient)
	}
}
