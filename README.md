Log Stream CLI Plugin
[![Concourse Badge][ci-badge]][ci-tests]
====================

The Log Stream CLI Plugin is a [CF CLI][cf-cli] plugin to retrieve logs from
a Loggregator V2 stream

### Installing Plugin

#### From CF-Community

```
cf install-plugin -r CF-Community "log-stream"
```

#### From Binary Release

1. Download the binary for the [latest release][latest-release] for your
   platform.
1. Install it into the cf cli:

```
cf install-plugin download/path/log-stream-cli
```

#### From Source Code

Make sure to have the [latest Go toolchain][golang-dl] installed.

```
go get code.cloudfoundry.org/log-stream-cli/cmd/log-stream-cli
cf install-plugin $GOPATH/bin/log-stream-cli
```

### Usage

#### Log Stream
```
$ cf log-stream --help
NAME:
   log-stream - Stream all messages of all types from Loggregator

USAGE:
   log-stream <source-id> [<source-id>] [options]

OPTIONS:
   --shard-id       Distribute logs between multiple consumers
   --type, -t       Filter the streamed logs. Available: 'log','event','counter','gauge','timer'. Allows multiple.

```

The `source-id` can either be the application name, the application guid or the name of the component (e.g. `doppler`, `uaa`, `gorouter`...). You can provide as many source-ids as you want. If you are a platform admin (you have the the `logs.admin` scope from UAA), then you can also omit the source-id and you see all messages available in the RLP.

[cf-cli]: https://code.cloudfoundry.org/cli
[ci-badge]: https://loggregator.ci.cf-app.com/api/v1/pipelines/products/jobs/log-stream-cli-tests/badge
[ci-tests]: https://loggregator.ci.cf-app.com/teams/main/pipelines/products/jobs/log-stream-cli-tests
[golang-dl]: https://golang.org/dl/
[latest-release]: https://github.com/cloudfoundry/log-stream-cli/releases/latest

