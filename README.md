# High-Performance-Log-Processor

## Summary

`High-Performance-Log-Processor` is a streaming data-processing project for parsing, filtering, and aggregating large volumes of log input efficiently. It is meant to emphasize throughput, allocation control, and pipeline design.

## Why This Project Exists

This project is meant to teach:

- how to structure high-throughput pipelines in Go,
- batching and backpressure tradeoffs,
- parsing unstructured or semi-structured data efficiently,
- where performance work actually matters in streaming systems.

## Planned Capabilities

- Read logs from files, stdin, or network input.
- Parse lines into structured records.
- Filter or transform records through configurable stages.
- Produce summary statistics or transformed output.

## Scope Boundaries

The initial goal for this project is to process logs after they have already been written rather than to build a live log collection system.

That means the first version should focus on:

- reading from existing log files,
- parsing lines into structured records,
- filtering, transforming, and aggregating records efficiently,
- producing summaries or transformed output.

Possible later expansions:

1. Support `stdin` so the processor can be used in shell pipelines.
2. Support long-running input sources such as `tail -f`.
3. Explore optional network ingestion only after the file and stdin pipeline is solid.

This keeps the learning focus on throughput, parsing, batching, buffering, and pipeline design instead of operational concerns like transport protocols, delivery guarantees, or full log shipping infrastructure.

## Architecture Sketch

- A reader stage ingests raw log lines.
- Parser workers transform lines into records.
- Downstream stages filter, aggregate, or export results.
- The design should make buffering and batching policies explicit.

## Milestones

1. Ingest lines and parse a simple log format.
2. Add a multi-stage concurrent processing pipeline.
3. Add aggregation and performance benchmarks.
4. Tune allocations, batching, and throughput under load.
5. Add `stdin` support for pipeline-style usage.
6. Evaluate whether network ingestion is still worthwhile after the file and stdin workflow feels complete.

## Current Status

This project is no longer just scaffolded.

Current progress includes:

- a tokenizer for the initial `key=value` log format with quoted values,
- parsing into `dashboard.LogRecord`,
- realistic sample fixtures in `testdata/`,
- shared path aggregation helpers,
- shared window bucketing helpers,
- implemented path and window metric functions for:
  - requests by window,
  - levels by window,
  - warn/error counts by window,
  - status classes by window,
  - status codes by window,
  - metrics by path,
  - metrics by path and window,
  - latency by path,
  - slow requests by path,
  - slow requests by window,
  - error rate by path,
  - error rate by window,
- JSON-ready dashboard structs,
- a human-readable aligned text renderer for default CLI output,
- a focused automated test suite for the implemented behavior.

The project is still in active development. Malformed-input handling is still being refined, the CLI/output surface is still minimal, and the parsing path still has some performance and robustness cleanup left.

## Development Notes

Current useful commands:

- `go run ./cmd/High-Performance-Log-Processor`
- `go build ./cmd/High-Performance-Log-Processor`
- `go test ./...`

Additional design notes live in `doc/design.md`.

## Usage

Current CLI behavior is intentionally small.

Supported modes today:

- parse a log file and choose a metric/report to emit
- parse `ping` output from `stdin`
- print `help`
- print `version`

Current report options:

- `MetricsByPath`
- `LatencyByPath`
- `SlowRequestsByPath`
- `ErrorRateByPath`
- `RequestsByWindow`
- `LevelsByWindow`
- `WarnAndErrorCountsByWindow`
- `StatusClassesByWindow`
- `StatusCodesByWindow`
- `MetricsByPathAndWindow`
- `SlowRequestsByWindow`
- `ErrorRateByWindow`

Current flags:

- `--time <duration>` for windowed reports
- `--output JSON`

Current CLI defaults:

- if no report is given for file input, the program defaults to `MetricsByPath`
- if no bucket is given for windowed reports, the program defaults to `5m`
- if `--output JSON` is not provided, the program prints aligned text output

Examples:

```bash
go run ./cmd/High-Performance-Log-Processor ./testdata/access.log
```

```bash
go run ./cmd/High-Performance-Log-Processor ./testdata/access.log MetricsByPath --output JSON
```

```bash
go run ./cmd/High-Performance-Log-Processor ./testdata/access.log RequestsByWindow --time 5m --output JSON
```

```bash
ping 8.8.8.8 | go run ./cmd/High-Performance-Log-Processor ping
```

## Example Output

The current file-processing path prints aligned text output by default.

Example:

```text
/api/login               req=     4 info=   1 warn=   3 err=   0 2xx=   1 4xx=   3 5xx=   0 avg_ms=  14 max_ms=  19
/api/products            req=     2 info=   2 warn=   0 err=   0 2xx=   2 4xx=   0 5xx=   0 avg_ms=  19 max_ms=  21
```

When `--output JSON` is selected, the current file-processing path prints one JSON object per result item.

Example:

```json
{"path":"/api/login","request_count":4,"level_counts":{"info_count":1,"warn_count":3},"status_counts":{"status_2xx":1,"status_4xx":3},"latency":{"count":4,"total_ms":59,"average_ms":14,"max_ms":19}}
{"path":"/api/products","request_count":2,"level_counts":{"info_count":2},"status_counts":{"status_2xx":2},"latency":{"count":2,"total_ms":39,"average_ms":19,"max_ms":21}}
```

The exact output depends on the input data and the currently selected report. If no report is given for file input, the built-in default is `MetricsByPath`.

## Project Structure

```text
cmd/High-Performance-Log-Processor/    processor entrypoint and parser helpers
internal/                              dashboard, aggregation, and CLI internals
pkg/                                   optional reusable processing components
doc/                                   design notes and benchmark ideas
scripts/                               helper scripts
testdata/                              sample log fixtures
```
