package dashboard

import "time"

// LogRecord represents one parsed log line ready for aggregation.
type LogRecord struct {
	TS         time.Time `json:"ts,omitempty"`
	Level      string    `json:"level,omitempty"`
	RequestID  string    `json:"request_id,omitempty"`
	Method     string    `json:"method,omitempty"`
	Path       string    `json:"path,omitempty"`
	Status     int       `json:"status,omitempty"`
	DurationMS int       `json:"duration_ms,omitempty"`
	Bytes      int       `json:"bytes,omitempty"`
	IP         string    `json:"ip,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	Message    string    `json:"message,omitempty"`
}

// TimeWindow describes the inclusive-exclusive period a metric bucket covers.
type TimeWindow struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

// RequestVolumePoint represents the total number of requests seen in one time bucket.
type RequestVolumePoint struct {
	Window       TimeWindow `json:"window,omitempty"`
	RequestCount int        `json:"request_count,omitempty"`
}

// LevelCounts contains totals split by log level.
type LevelCounts struct {
	InfoCount  int `json:"info_count,omitempty"`
	WarnCount  int `json:"warn_count,omitempty"`
	ErrorCount int `json:"error_count,omitempty"`
}

// LevelVolumePoint represents log level counts for one time bucket.
type LevelVolumePoint struct {
	Window TimeWindow  `json:"window,omitempty"`
	Counts LevelCounts `json:"counts,omitempty"`
}

// StatusClassCounts contains totals for each HTTP status class.
type StatusClassCounts struct {
	Status1xx int `json:"status_1xx,omitempty"`
	Status2xx int `json:"status_2xx,omitempty"`
	Status3xx int `json:"status_3xx,omitempty"`
	Status4xx int `json:"status_4xx,omitempty"`
	Status5xx int `json:"status_5xx,omitempty"`
}

// StatusClassVolumePoint represents status-class totals for one time bucket.
type StatusClassVolumePoint struct {
	Window TimeWindow        `json:"window,omitempty"`
	Counts StatusClassCounts `json:"counts,omitempty"`
}

// StatusCodeCount represents the number of times a specific HTTP status appears.
type StatusCodeCount struct {
	StatusCode int `json:"status_code,omitempty"`
	Count      int `json:"count,omitempty"`
}

// StatusCodeVolumePoint represents exact status-code totals for one time bucket.
// Example output could include counts for 200, 201, 400, 404, 429, 500, and 503.
type StatusCodeVolumePoint struct {
	Window TimeWindow        `json:"window,omitempty"`
	Counts []StatusCodeCount `json:"counts,omitempty"`
}

// LatencySummary represents basic latency statistics plus slow-request thresholds.
type LatencySummary struct {
	Count         int `json:"count,omitempty"`
	TotalMs       int `json:"total_ms,omitempty"`
	AverageMS     int `json:"average_ms,omitempty"`
	MaxMS         int `json:"max_ms,omitempty"`
	SlowOver100MS int `json:"slow_over_100_ms,omitempty"`
	SlowOver250MS int `json:"slow_over_250_ms,omitempty"`
	SlowOver500MS int `json:"slow_over_500_ms,omitempty"`
}

// PathMetrics represents aggregated metrics for a specific API path or route.
// It should capture request volume, level totals, status-class totals, and latency data.
type PathMetrics struct {
	Path         string            `json:"path,omitempty"`
	RequestCount int               `json:"request_count,omitempty"`
	LevelCounts  LevelCounts       `json:"level_counts,omitempty"`
	StatusCounts StatusClassCounts `json:"status_counts,omitempty"`
	Latency      LatencySummary    `json:"latency,omitempty"`
}

// PathWindowMetrics represents per-path metrics within a single time bucket.
// This is useful for answering questions like:
// "What happened on /api/products between 09:00 and 09:05?"
type PathWindowMetrics struct {
	Window TimeWindow    `json:"window,omitempty"`
	Paths  []PathMetrics `json:"paths,omitempty"`
}

// WindowBucket represents one time bucket plus the raw records that belong to it.
// Windowed metric functions can reuse this type and then decide for themselves
// how to aggregate the records inside each bucket.
type WindowBucket struct {
	Window  TimeWindow  `json:"window,omitempty"`
	Records []LogRecord `json:"records,omitempty"`
}

// PathLatencyMetrics represents latency-only summaries for a specific path.
type PathLatencyMetrics struct {
	Path    string         `json:"path,omitempty"`
	Latency LatencySummary `json:"latency,omitempty"`
}

// BucketSize defines the aggregation interval, for example:
// 1 minute, 5 minutes, or 1 hour.
type BucketSize time.Duration
