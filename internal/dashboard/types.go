package dashboard

import "time"

// LogRecord represents one parsed log line ready for aggregation.
type LogRecord struct {
	TS         time.Time `json:"ts"`
	Level      string    `json:"level"`
	RequestID  string    `json:"request_id"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Status     int       `json:"status"`
	DurationMS int       `json:"duration_ms"`
	Bytes      int       `json:"bytes"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	Message    string    `json:"message"`
}

// TimeWindow describes the inclusive-exclusive period a metric bucket covers.
type TimeWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// RequestVolumePoint represents the total number of requests seen in one time bucket.
type RequestVolumePoint struct {
	Window       TimeWindow `json:"window"`
	RequestCount int        `json:"request_count"`
}

// LevelCounts contains totals split by log level.
type LevelCounts struct {
	InfoCount  int `json:"info_count"`
	WarnCount  int `json:"warn_count"`
	ErrorCount int `json:"error_count"`
}

// LevelVolumePoint represents log level counts for one time bucket.
type LevelVolumePoint struct {
	Window TimeWindow  `json:"window"`
	Counts LevelCounts `json:"counts"`
}

// StatusClassCounts contains totals for each HTTP status class.
type StatusClassCounts struct {
	Status1xx int `json:"status_1xx"`
	Status2xx int `json:"status_2xx"`
	Status3xx int `json:"status_3xx"`
	Status4xx int `json:"status_4xx"`
	Status5xx int `json:"status_5xx"`
}

// StatusClassVolumePoint represents status-class totals for one time bucket.
type StatusClassVolumePoint struct {
	Window TimeWindow       `json:"window"`
	Counts StatusClassCounts `json:"counts"`
}

// StatusCodeCount represents the number of times a specific HTTP status appears.
type StatusCodeCount struct {
	StatusCode int `json:"status_code"`
	Count      int `json:"count"`
}

// StatusCodeVolumePoint represents exact status-code totals for one time bucket.
// Example output could include counts for 200, 201, 400, 404, 429, 500, and 503.
type StatusCodeVolumePoint struct {
	Window TimeWindow        `json:"window"`
	Counts []StatusCodeCount `json:"counts"`
}

// LatencySummary represents basic latency statistics plus slow-request thresholds.
type LatencySummary struct {
	Count         int `json:"count"`
	TotalMs       int `json:"total_ms"`
	AverageMS     int `json:"average_ms"`
	MaxMS         int `json:"max_ms"`
	SlowOver100MS int `json:"slow_over_100_ms"`
	SlowOver250MS int `json:"slow_over_250_ms"`
	SlowOver500MS int `json:"slow_over_500_ms"`
}

// PathMetrics represents aggregated metrics for a specific API path or route.
// It should capture request volume, level totals, status-class totals, and latency data.
type PathMetrics struct {
	Path         string            `json:"path"`
	RequestCount int               `json:"request_count"`
	LevelCounts  LevelCounts       `json:"level_counts"`
	StatusCounts StatusClassCounts `json:"status_counts"`
	Latency      LatencySummary    `json:"latency"`
}

// PathWindowMetrics represents per-path metrics within a single time bucket.
// This is useful for answering questions like:
// "What happened on /api/products between 09:00 and 09:05?"
type PathWindowMetrics struct {
	Window TimeWindow    `json:"window"`
	Paths  []PathMetrics `json:"paths"`
}

// WindowBucket represents one time bucket plus the raw records that belong to it.
// Windowed metric functions can reuse this type and then decide for themselves
// how to aggregate the records inside each bucket.
type WindowBucket struct {
	Window  TimeWindow   `json:"window"`
	Records []LogRecord `json:"records"`
}

// PathLatencyMetrics represents latency-only summaries for a specific path.
type PathLatencyMetrics struct {
	Path    string         `json:"path"`
	Latency LatencySummary `json:"latency"`
}

// BucketSize defines the aggregation interval, for example:
// 1 minute, 5 minutes, or 1 hour.
type BucketSize time.Duration
