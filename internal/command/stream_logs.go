package command

import (
	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"context"
	"encoding/json"
	"github.com/cloudfoundry/log-stream-cli/internal/presentation"
	"github.com/cloudfoundry/log-stream-cli/internal/rlp"
	"github.com/gogo/protobuf/jsonpb"
	"io"
	"log"
	"net/http"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func StreamLogs(logStreamUrl string, client doer, writer io.Writer, options ...ApplyOptionFn) {
	c := loggregator.NewRLPGatewayClient(
		logStreamUrl,
		loggregator.WithRLPGatewayClientLogger(log.New(NewDedupeWriter(writer), "", 0)),
		loggregator.WithRLPGatewayHTTPClient(client),
	)

	opts := &streamLogsOptions{}
	for _, apply := range options {
		apply(opts)
	}

	r, err := rlp.MakeRequest(opts.sourceIDs, opts.metricTypes, opts.shardID)
	if err != nil {
		log.Fatal(err)
	}

	streamLogs(c, r, writer)
}

func streamLogs(c *loggregator.RLPGatewayClient, r *loggregator_v2.EgressBatchRequest, writer io.Writer) {
	es := c.Stream(context.Background(), r)
	marshaler := jsonpb.Marshaler{}
	for {
		for _, e := range es() {
			switch e.Message.(type) {
			case *loggregator_v2.Envelope_Log:
				bytes, err := json.Marshal(presentation.BuildBase64DecodedLog(e))
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

type ApplyOptionFn func(*streamLogsOptions)

type streamLogsOptions struct {
	sourceIDs   []string
	metricTypes []string
	shardID     string
}

func WithSourceIDs(sourceIDs []string) ApplyOptionFn {
	return func(opt *streamLogsOptions) {
		opt.sourceIDs = sourceIDs
	}
}

func WithMetricTypes(metricTypes []string) ApplyOptionFn {
	return func(opt *streamLogsOptions) {
		opt.metricTypes = metricTypes
	}
}

func WithShardID(shardID string) ApplyOptionFn {
	return func(opt *streamLogsOptions) {
		opt.shardID = shardID
	}
}
