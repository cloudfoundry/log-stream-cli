package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	"github.com/cloudfoundry/log-stream-cli/internal/log_stream_plugin"

	flags "github.com/jessevdk/go-flags"
)

type CFLogStreamCLI struct{}

type CLIFlags struct {
	MetricTypes []string `short:"t" long:"type"`
}

func (c CFLogStreamCLI) Run(conn plugin.CliConnection, args []string) {
	accessToken, err := conn.AccessToken()
	if err != nil {
		log.Fatal("unable to retrieve access token", err)
	}

	skipSSL, err := conn.IsSSLDisabled()
	if err != nil {
		log.Fatal("unable to retrieve skip ssl flag", err)
	}

	logStreamEndpoint, err := logStreamEndpoint(conn)
	if err != nil {
		log.Fatal("invalid log stream endpoint", err)
	}

	var cliFlags CLIFlags
	parser := flags.NewParser(&cliFlags, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Fatal("error parsing flags", err)
	}

	if args, err = parser.ParseArgs(args); err != nil {
		log.Fatal("error parsing args", err)
	}

	switch args[0] {
	case "log-stream":
		command.StreamLogs(
			logStreamEndpoint,
			log_stream_plugin.NewDoer(accessToken, skipSSL),
			os.Stdout,
			command.WithSourceIDs(args[1:]),
			command.WithMetricTypes(cliFlags.MetricTypes),
		)
	}

	return
}

// version is set via ldflags at compile time. It should be JSON encoded
// plugin.VersionType. If it does not unmarshal, the plugin version will be
// left empty.
var version string

func (c CFLogStreamCLI) GetMetadata() plugin.PluginMetadata {
	var v plugin.VersionType
	// Ignore the error. If this doesn't unmarshal, then we want the default
	// VersiionType.
	_ = json.Unmarshal([]byte(version), &v)

	return plugin.PluginMetadata{
		Name:    "log-stream",
		Version: v,
		Commands: []plugin.Command{
			{
				Name:     "log-stream",
				HelpText: "Stream all messages of all types from Loggregator",
				UsageDetails: plugin.Usage{
					Usage: "log-stream <source-id> [<source-id>] [options]",
					Options: map[string]string{
						"-type, -t": "Filter the streamed logs. Available: 'log','event','counter','gauge','timer'. Allows multiple.",
					},
				},
			},
		},
	}
}

func logStreamEndpoint(conn plugin.CliConnection) (string, error) {
	apiEndpoint, err := conn.ApiEndpoint()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("://(api)")
	return re.ReplaceAllString(apiEndpoint, "://log-stream"), nil
}

func main() {
	plugin.Start(CFLogStreamCLI{})
}
