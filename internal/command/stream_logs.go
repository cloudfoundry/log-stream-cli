package command

import (
	"context"
	"io"
	"log"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/cloudfoundry/log-stream-cli/internal/log_stream_plugin"
	"github.com/gogo/protobuf/jsonpb"
)

func StreamLogs(logStreamUrl string, doer log_stream_plugin.Doer, writer io.Writer) {
	c := loggregator.NewRLPGatewayClient(
		logStreamUrl,
		loggregator.WithRLPGatewayClientLogger(log.New(log_stream_plugin.NewDedupeWriter(writer), "", 0)),
		loggregator.WithRLPGatewayHTTPClient(doer),
	)

	es := c.Stream(context.Background(), selectors)

	marshaler := jsonpb.Marshaler{}
	for {
		for _, e := range es() {
			if err := marshaler.Marshal(writer, e); err != nil {
				log.Fatal(err)
			}
			writer.Write([]byte("\n"))
		}
	}
}

var selectors = &loggregator_v2.EgressBatchRequest{
	Selectors: []*loggregator_v2.Selector{
		{
			Message: &loggregator_v2.Selector_Log{
				Log: &loggregator_v2.LogSelector{},
			},
		},
		{
			Message: &loggregator_v2.Selector_Counter{
				Counter: &loggregator_v2.CounterSelector{},
			},
		},
		{
			Message: &loggregator_v2.Selector_Event{
				Event: &loggregator_v2.EventSelector{},
			},
		},
		{
			Message: &loggregator_v2.Selector_Gauge{
				Gauge: &loggregator_v2.GaugeSelector{},
			},
		},
		{
			Message: &loggregator_v2.Selector_Timer{
				Timer: &loggregator_v2.TimerSelector{},
			},
		},
	},
}
