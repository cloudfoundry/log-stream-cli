package command

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"context"
	"encoding/json"
	"github.com/cloudfoundry/log-stream-cli/internal/presentation"
	"github.com/cloudfoundry/log-stream-cli/internal/rlp"
	"io"
	"log"
	"net/http"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type appProvider interface {
	GetApps() ([]plugin_models.GetAppsModel, error)
}

func StreamLogs(logStreamUrl string, client doer, ap appProvider, writer io.Writer, options ...ApplyOptionFn) {
	c := loggregator.NewRLPGatewayClient(
		logStreamUrl,
		loggregator.WithRLPGatewayClientLogger(log.New(NewDedupeWriter(writer), "", 0)),
		loggregator.WithRLPGatewayHTTPClient(client),
	)

	opts := &streamLogsOptions{}
	for _, apply := range options {
		apply(opts)
	}

	sourceIDs := replaceAppNamesWithGuids(opts.sourceIDs, ap)
	r, err := rlp.MakeRequest(sourceIDs, opts.metricTypes, opts.shardID)
	if err != nil {
		log.Fatal(err)
	}

	streamLogs(c, r, writer)
}

func replaceAppNamesWithGuids(sids []string, ap appProvider) []string {
	apps, err := ap.GetApps()
	if err != nil {
		log.Printf("Warning, unable to retrieve apps: %s. Using raw input", err)
		return sids
	}
	for _, a := range apps {
		for i, s := range sids {
			if a.Name == s {
				sids[i] = a.Guid
			}
		}
	}
	return sids
}

func streamLogs(c *loggregator.RLPGatewayClient, r *loggregator_v2.EgressBatchRequest, writer io.Writer) {
	es := c.Stream(context.Background(), r)
	for {
		for _, e := range es() {
			presEnv, err := presentation.Envelope(e)
			if err != nil {
				log.Fatal(err)
			}

			envString, err := json.Marshal(presEnv); if err != nil {
				log.Fatal(err)
			}

			writer.Write(envString)
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
