package dashboard

import "time"

// LogRecord represents one parsed log line ready for aggregation.
type LogRecord struct {
	TS         time.Time
	Level      string
	RequestID  string
	Method     string
	Path       string
	Status     int
	DurationMS int
	Bytes      int
	IP         string
	UserAgent  string
	Message    string
}

// TimeWindow describes the inclusive-exclusive period a metric bucket covers.
type TimeWindow struct {
	Start time.Time
	End   time.Time
}

// RequestVolumePoint represents the total number of requests seen in one time bucket.
type RequestVolumePoint struct {
	Window       TimeWindow
	RequestCount int
}

// LevelCounts contains totals split by log level.
type LevelCounts struct {
	InfoCount  int
	WarnCount  int
	ErrorCount int
}

// LevelVolumePoint represents log level counts for one time bucket.
type LevelVolumePoint struct {
	Window TimeWindow
	Counts LevelCounts
}

// StatusClassCounts contains totals for each HTTP status class.
type StatusClassCounts struct {
	Status1xx int
	Status2xx int
	Status3xx int
	Status4xx int
	Status5xx int
}

// StatusClassVolumePoint represents status-class totals for one time bucket.
type StatusClassVolumePoint struct {
	Window TimeWindow
	Counts StatusClassCounts
}

// StatusCodeCount represents the number of times a specific HTTP status appears.
type StatusCodeCount struct {
	StatusCode int
	Count      int
}

// StatusCodeVolumePoint represents exact status-code totals for one time bucket.
// Example output could include counts for 200, 201, 400, 404, 429, 500, and 503.
type StatusCodeVolumePoint struct {
	Window TimeWindow
	Counts []StatusCodeCount
}

// LatencySummary represents basic latency statistics plus slow-request thresholds.
type LatencySummary struct {
	Count         int
	TotalMs       int
	AverageMS     int
	MaxMS         int
	SlowOver100MS int
	SlowOver250MS int
	SlowOver500MS int
}

// PathMetrics represents aggregated metrics for a specific API path or route.
// It should capture request volume, level totals, status-class totals, and latency data.
type PathMetrics struct {
	Path         string
	RequestCount int
	LevelCounts  LevelCounts
	StatusCounts StatusClassCounts
	Latency      LatencySummary
}

// PathWindowMetrics represents per-path metrics within a single time bucket.
// This is useful for answering questions like:
// "What happened on /api/products between 09:00 and 09:05?"
type PathWindowMetrics struct {
	Window TimeWindow
	Paths  []PathMetrics
}

// WindowBucket represents one time bucket plus the raw records that belong to it.
// Windowed metric functions can reuse this type and then decide for themselves
// how to aggregate the records inside each bucket.
type WindowBucket struct {
	Window  TimeWindow
	Records []LogRecord
}

// PathLatencyMetrics represents latency-only summaries for a specific path.
type PathLatencyMetrics struct {
	Path    string
	Latency LatencySummary
}

// BucketSize defines the aggregation interval, for example:
// 1 minute, 5 minutes, or 1 hour.
type BucketSize time.Duration
