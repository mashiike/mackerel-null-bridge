# mackerel-null-bridge

![Latest GitHub release](https://img.shields.io/github/release/mashiike/mackerel-null-bridge.svg)
![Github Actions test](https://github.com/mashiike/mackerel-null-bridge/workflows/Test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/mashiike/mackerel-null-bridge)](https://goreportcard.com/report/mashiike/mackerel-null-bridge) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/mashiike/mackerel-null-bridge/blob/master/LICENSE)

A command line tool for filling missing metric values on Mackerel.

## Description

When sending error metrics, etc., you may be forced to send them intermittently.
When monitoring such intermittent metrics in Mackerel, alerts may not close automatically.
This tool is designed to be run periodically, so it will fetch the values of the specified metrics not more than 15 minutes apart and interpolate the missing values. As a result, it expects to close alerts automatically, even for intermittent metrics.

Note: As of v0.0.0, only service metrics are still supported. This is because I can't think of any case where intermittent values are sent in the host metric. 
## Install

### binary packages

[Releases](https://github.com/mashiike/mackerel-null-bridge/releases).

### Homebrew tap

```console
$ brew install mashiike/tap/mackerel-null-bridge
```
## Usage

### as CLI command

```console
$mackerel-null-bridge
NAME:
   mackerel-null-bridge - A command line tool for filling missing metric values on Mackerel.

USAGE:
   mackerel-null-bridge --config <config file> --apikey <Mackerel APIKEY>

VERSION:
   0.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --apikey value, -k value  for access mackerel API (default: *********) [$MACKEREL_APIKEY]
   --config value, -c value  config file path, can set multiple [$CONFIG_FILE]
   --deploy                  deploy flag (cli only) (default: false)
   --dry-run                 dry-run flag (lambda only) (default: false) [$DRY_RUN]
   --log-level value         output log level (default: info) [$LOG_LEVEL]
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

### as AWS Lambda function

`mackerel-null-bridge` binary also runs as AWS Lambda function. 
mackerel-null-bridge implicitly behaves as a run command when run as a bootstrap with a Lambda Function


CLI options can be specified from environment variables. For example, when `MACKEREL_APIKEY` environment variable is set, the value is set to `-apikey` option.

Example Lambda functions configuration.

```json
{
  "FunctionName": "mackerel-null-bridge",
  "Environment": {
    "Variables": {
      "CONFIG_FILE": "config.yaml",
      "MACKEREL_APIKEY": "<Mackerel API KEY>"
    }
  },
  "Handler": "bootstrap",
  "MemorySize": 128,
  "Role": "arn:aws:iam::0123456789012:role/lambda-function",
  "Runtime": "provided.al2",
  "Timeout": 300
}
```

### Configuration file

YAML format.

```yaml
required_version: ">=0.0.0"

targets:
  - service: prod
    metric_name: hoge.fuga.piyo
    value: 0.0
    delay_seconds: 300
```

# Special Thanks
[@handlename](https://github.com/handlename) gave me the idea to name this tool.

## LICENSE

MIT
