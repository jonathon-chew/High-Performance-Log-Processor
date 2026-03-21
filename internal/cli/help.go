package cli

const usageText string = `USAGE:
  High-Performance-Log-Processor <file> [metric] [--time <duration>] [--output <format>]
  High-Performance-Log-Processor ping
  High-Performance-Log-Processor help
  High-Performance-Log-Processor version

CURRENT BEHAVIOUR:
  - If no metric is provided for file input, the program defaults to MetricsByPath.
  - Windowed metrics use a default bucket size of 5m when --time is not provided.
  - Output defaults to plain text unless --output JSON is provided.
  - ping returns immediately into ping mode and ignores any arguments that follow it.

METRICS:
  MetricsByPath
  LatencyByPath
  SlowRequestsByPath
  ErrorRateByPath
  RequestsByWindow
  LevelsByWindow
  WarnAndErrorCountsByWindow
  StatusClassesByWindow
  StatusCodesByWindow
  MetricsByPathAndWindow
  SlowRequestsByWindow
  ErrorRateByWindow

FLAGS:
  --time <duration>     Bucket size for windowed metrics, for example 1m, 5m, 1h
  --output <format>     Currently supports JSON

EXAMPLES:
  go run ./cmd/High-Performance-Log-Processor ./testdata/access.log
  go run ./cmd/High-Performance-Log-Processor ./testdata/access.log MetricsByPath
  go run ./cmd/High-Performance-Log-Processor ./testdata/access.log RequestsByWindow --time 5m --output JSON
  ping 8.8.8.8 | go run ./cmd/High-Performance-Log-Processor ping
`
