package main

import (
	"log"
	"os"
	"regexp"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	"github.com/cloudfoundry/log-stream-cli/internal/log_stream_plugin"
)

type CFLogStreamCLI struct{}

const UNINSTALL_FLAG = "CLI-MESSAGE-UNINSTALL"

func (c CFLogStreamCLI) Run(conn plugin.CliConnection, args []string) {
	if args[0] == UNINSTALL_FLAG {
		return
	}
	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	ui := terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

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

	command.StreamLogs(logStreamEndpoint, log_stream_plugin.NewDoer(accessToken, skipSSL), ui)

	return
}

func (c CFLogStreamCLI) GetMetadata() plugin.PluginMetadata {
	v := plugin.VersionType{
		Major: 0,
		Minor: 0,
		Build: 1,
	}

	return plugin.PluginMetadata{
		Name:    "log-stream",
		Version: v,
		Commands: []plugin.Command{
			{
				Name:     "log-stream",
				HelpText: "Stream all messages of all types from Loggregator",
				UsageDetails: plugin.Usage{
					Usage: "log-stream",
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
