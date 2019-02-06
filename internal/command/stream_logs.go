package command

import (
	"context"
	"encoding/json"
	"io"
	"log"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/cloudfoundry/log-stream-cli/internal/log_stream_plugin"
	"github.com/gogo/protobuf/jsonpb"
)

func StreamLogs(sourceIDs []string, logStreamUrl string, doer log_stream_plugin.Doer, writer io.Writer) {
	c := loggregator.NewRLPGatewayClient(
		logStreamUrl,
		loggregator.WithRLPGatewayClientLogger(log.New(log_stream_plugin.NewDedupeWriter(writer), "", 0)),
		loggregator.WithRLPGatewayHTTPClient(doer),
	)

	es := c.Stream(context.Background(), req(sourceIDs))

	marshaler := jsonpb.Marshaler{}
	for {
		for _, e := range es() {
			switch e.Message.(type) {
			case *loggregator_v2.Envelope_Log:
				bytes, err := json.Marshal(log_stream_plugin.BuildBase64DecodedLog(e))
				if err != nil {
					log.Fatal("error marshalling", err)
				}

				_, err = writer.Write(bytes)
				if err != nil {
					log.Fatal(err)
				}
			default:
				if err := marshaler.Marshal(writer, e); err != nil {
					log.Fatal(err)
				}
			}
			writer.Write([]byte("\n"))
		}
	}
}

func req(sourceIDs []string) *loggregator_v2.EgressBatchRequest {
	var s []*loggregator_v2.Selector
	for _, sourceId := range sourceIDs {
		s = append(s, selectors(sourceId)...)
	}

	if len(s) == 0 {
		s = selectors("")
	}

	return &loggregator_v2.EgressBatchRequest{
		Selectors: s,
	}
}

func selectors(sourceId string) []*loggregator_v2.Selector {
	return []*loggregator_v2.Selector{
		{
			SourceId: sourceId,
			Message: &loggregator_v2.Selector_Log{
				Log: &loggregator_v2.LogSelector{},
			},
		},
		{
			SourceId: sourceId,
			Message: &loggregator_v2.Selector_Counter{
				Counter: &loggregator_v2.CounterSelector{},
			},
		},
		{
			SourceId: sourceId,
			Message: &loggregator_v2.Selector_Event{
				Event: &loggregator_v2.EventSelector{},
			},
		},
		{
			SourceId: sourceId,
			Message: &loggregator_v2.Selector_Gauge{
				Gauge: &loggregator_v2.GaugeSelector{},
			},
		},
		{
			SourceId: sourceId,
			Message: &loggregator_v2.Selector_Timer{
				Timer: &loggregator_v2.TimerSelector{},
			},
		},
	}
}
