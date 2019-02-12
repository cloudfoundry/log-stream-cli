package command

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

var metricTypeToSelector = map[string]loggregator_v2.Selector{
	"log": {
		Message: &loggregator_v2.Selector_Log{
			Log: &loggregator_v2.LogSelector{},
		},
	},
	"counter": {
		Message: &loggregator_v2.Selector_Counter{
			Counter: &loggregator_v2.CounterSelector{},
		},
	},
	"event": {
		Message: &loggregator_v2.Selector_Event{
			Event: &loggregator_v2.EventSelector{},
		},
	},
	"gauge": {
		Message: &loggregator_v2.Selector_Gauge{
			Gauge: &loggregator_v2.GaugeSelector{},
		},
	},
	"timer": {
		Message: &loggregator_v2.Selector_Timer{
			Timer: &loggregator_v2.TimerSelector{},
		},
	},
}

func MakeRequest(sourceIDs []string, metricTypes []string, shardID string) (*loggregator_v2.EgressBatchRequest, error) {
	var invalid []string
	for _, t := range metricTypes {
		if _, ok := metricTypeToSelector[t]; !ok {
			invalid = append(invalid, t)
		}
	}

	if len(invalid) > 0 {
		return nil, fmt.Errorf("invalid metric type(s): %s", strings.Join(invalid, ", "))
	}

	if len(sourceIDs) == 0 {
		sourceIDs = []string{""}
	}

	if len(metricTypes) == 0 {
		metricTypes = []string{"log", "counter", "event", "gauge", "timer"}
	}

	return &loggregator_v2.EgressBatchRequest{
		ShardId: shardID,
		Selectors: crossProduct(
			sourceIDs,
			metricTypes,
		),
	}, nil
}

func getSelectorOfType(metricType string, sourceID string) *loggregator_v2.Selector {
	selector, ok := metricTypeToSelector[strings.ToLower(metricType)]
	if !ok {
		panic("error: invalid metric")
	}
	selector.SourceId = sourceID
	return &selector
}

func crossProduct(sourceIDs []string, messageTypes []string) []*loggregator_v2.Selector {
	var ss []*loggregator_v2.Selector

	for _, sid := range sourceIDs {
		for _, mt := range messageTypes {
			ss = append(
				ss,
				getSelectorOfType(mt, sid),
			)
		}
	}

	return ss
}
