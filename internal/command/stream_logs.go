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

func StreamLogs(logStreamUrl string, doer log_stream_plugin.Doer, writer io.Writer, options ...applyOptionFn) {
	c := loggregator.NewRLPGatewayClient(
		logStreamUrl,
		loggregator.WithRLPGatewayClientLogger(log.New(log_stream_plugin.NewDedupeWriter(writer), "", 0)),
		loggregator.WithRLPGatewayHTTPClient(doer),
	)

	opts := &streamLogsOptions{}
	for _, apply := range options {
		apply(opts)
	}

	r, err := log_stream_plugin.MakeRequest(opts.sourceIDs, opts.metricTypes, opts.shardID)
	if err != nil {
		log.Fatal(err)
	}
	es := c.Stream(context.Background(), r)

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

type applyOptionFn func(*streamLogsOptions)

type streamLogsOptions struct {
	sourceIDs   []string
	metricTypes []string
	shardID     string
}

func WithSourceIDs(sourceIDs []string) applyOptionFn {
	return func(opt *streamLogsOptions) {
		opt.sourceIDs = sourceIDs
	}
}

func WithMetricTypes(metricTypes []string) applyOptionFn {
	return func(opt *streamLogsOptions) {
		opt.metricTypes = metricTypes
	}
}

func WithShardID(shardID string) applyOptionFn {
	return func(opt *streamLogsOptions) {
		opt.shardID = shardID
	}
}
