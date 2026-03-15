# High-Performance-Log-Processor Design Notes

## Project Boundary

The initial version of this project is intended to process logs after they have already been written.

This project is not currently intended to be:

- a live traffic collector,
- a log shipping system,
- a networked observability agent.

The first target is file-based analysis, with possible later support for:

- `stdin`,
- long-running piped input such as `tail -f`,
- optional network ingestion after the file pipeline is solid.

## Input Format

The first log format is a plain-text `key=value` format with quoted values for fields that may contain spaces.

Example:

```text
ts=2026-03-14T09:00:00.112Z level=INFO req_id=0f4c9f01 method=GET path=/ status=200 duration_ms=8 bytes=5120 ip=203.0.113.10 ua="Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4)" msg="request complete"
```

Important fields:

- `ts`
- `level`
- `req_id`
- `method`
- `path`
- `status`
- `duration_ms`
- `bytes`
- `ip`
- `ua`
- `msg`

## Parsing Decisions

### Tokenization

The line format cannot be split with a naive space split because quoted values contain spaces.

The tokenizer design is:

- scan the line byte-by-byte,
- track whether parsing is currently inside double quotes,
- split on spaces only when not inside quotes,
- append the final token after the scan completes.

This keeps tokenization separate from field interpretation.

### Field Parsing

Field parsing is treated as a second stage after tokenization.

The intended parsing flow is:

1. Tokenize the line into `key=value` parts.
2. Split each token once on the first `=`.
3. Strip surrounding quotes from quoted values during value handling, not during tokenization.

This separation keeps the tokenizer small and makes later format changes easier to manage.

### Time Parsing

Timestamps use RFC3339-style values such as:

```text
2026-03-14T09:01:20.006Z
```

The intended Go parsing approach is:

- use `time.Parse`
- use `time.RFC3339Nano` or an equivalent Go layout

## Sample Data Strategy

Sample logs were created to support realistic throughput and aggregation work.

Current fixtures:

- `testdata/access.log`
- `testdata/access-large.log`

The large fixture is intended for throughput and aggregation testing.

The current dataset includes:

- normal application traffic,
- health checks,
- 4xx and 5xx responses,
- rate limiting,
- slow endpoints,
- repeated hot routes,
- bot-style probes.

Malformed-line scenarios are intentionally deferred to a later stage.

## Aggregation Goals

The main aggregation focus is:

- requests by time window,
- warnings and errors by time window,
- status-class counts by time window,
- status-code counts by time window,
- metrics by path,
- metrics by path and window,
- latency by path,
- slow requests by path,
- error rate by path,
- error rate by window.

Recommended default bucket sizes:

- `1m` for high detail,
- `5m` as the default operational view,
- `1h` for summary views.

## Metrics Design

### Shared Aggregation

`MetricsByPath` is treated as the core reusable path-level aggregation.

That aggregation is implemented through a shared helper:

- `aggregatePathMetrics(records []LogRecord) []PathMetrics`

This helper is intended to be the single source of truth for:

- request counts,
- level counts,
- status-class counts,
- latency totals,
- latency averages,
- maximum latency,
- slow-request thresholds.

### Windowing

Windowed functions should not duplicate path aggregation logic.

Instead, a separate helper groups raw records into time buckets:

- `groupRecordsByWindow(records []LogRecord, bucketSize BucketSize) []WindowBucket`

`WindowBucket` contains:

- the time window,
- the raw records in that window.

This helper is intentionally lower-level than `PathWindowMetrics` so it can be reused by multiple windowed functions, including:

- request counts by window,
- level totals by window,
- status totals by window,
- path metrics by window.

### Rolling Window Semantics

Current windowing is intentionally record-anchored rather than wall-clock aligned.

That means:

- a bucket starts at the first record in that bucket,
- subsequent records stay in that bucket until they exceed `bucketSize`,
- the next bucket starts at the next record that falls outside the previous bucket.

This is different from fixed clock windows such as:

- `09:00-09:05`
- `09:05-09:10`

The current function signatures permit this record-anchored interpretation.

`TimeWindow` is currently treated as the observed range of records in the bucket:

- `Start` is the first record timestamp in the bucket,
- `End` is the last record timestamp in the bucket.

## Latency Design

Average latency is not computed directly from the final output fields alone.

During aggregation, latency requires intermediate state:

- count,
- total duration,
- max duration,
- slow-over-threshold counts.

The average is computed after aggregation as:

```text
average = total duration / request count
```

This is why `LatencySummary` includes:

- `Count`
- `TotalMs`
- `AverageMS`
- `MaxMS`
- `SlowOver100MS`
- `SlowOver250MS`
- `SlowOver500MS`

Slow-threshold counts are currently cumulative:

- `>= 500ms` also contributes to `>= 250ms` and `>= 100ms`
- `>= 250ms` also contributes to `>= 100ms`

## Function-Level Decisions

### MetricsByPath

This function is intended to:

- aggregate all records by path,
- sort results by path before returning.

### MetricsByPathAndWindow

This function is intended to:

- group records into time buckets,
- aggregate each bucket by path using `aggregatePathMetrics`,
- return one `PathWindowMetrics` per bucket.

### ErrorRateByPath

This function returns `[]PathMetrics` rather than a narrower error-only type.

That design is intentional because `PathMetrics` already contains:

- per-path request totals,
- status-class totals.

This provides enough information for callers to calculate path-level error rates without introducing another dedicated return type.

### ErrorRateByWindow

This function is intended to:

- group records into windows,
- aggregate each window by path,
- sum status-class counts across all paths in the bucket,
- return one `StatusClassVolumePoint` per bucket.

It should not return one entry per path.

### LatencyByPath

This function is intended to:

- reuse `aggregatePathMetrics`,
- project each path down to `PathLatencyMetrics`.

### SlowRequestsByPath

This function is intended to return only paths that have at least one slow request.

The current interpretation of "slow issue" is:

- any path with a threshold count greater than zero.

The threshold checks are kept explicit for all three slow counters so the function remains correct even if threshold semantics change later.

## Organization Decisions

The dashboard package is split by concern rather than by one-struct-per-file.

Current file layout:

- `types.go`
- `path_metrics.go`
- `window_helpers.go`
- `window_metrics.go`

This keeps:

- shared types together,
- path aggregation logic together,
- window bucketing logic together,
- window-facing public functions together.

## Go Style Notes

Several style decisions were chosen intentionally:

- use maps for grouping by path,
- use small focused helpers instead of generic abstractions,
- prefer clear domain-specific helpers over reflection or over-generalized utilities,
- split package files by concern for readability.

Future repeated field-by-field accumulations may justify small helpers such as:

- adding one `StatusClassCounts` into another,
- adding one `LevelCounts` into another,
- updating latency accumulators from records.

These should remain narrow and purpose-specific.
