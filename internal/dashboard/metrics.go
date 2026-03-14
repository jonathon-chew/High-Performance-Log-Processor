package dashboard

import (
	"slices"
	"strings"
	"time"
)

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

// PathLatencyMetrics represents latency-only summaries for a specific path.
type PathLatencyMetrics struct {
	Path    string
	Latency LatencySummary
}

// BucketSize defines the aggregation interval, for example:
// 1 minute, 5 minutes, or 1 hour.
type BucketSize time.Duration

// RequestsByWindow should return total request counts for each time bucket,
// regardless of method, path, level, or status code.
func RequestsByWindow(records []LogRecord, bucketSize BucketSize) []RequestVolumePoint {
	panic("not implemented")
}

// LevelsByWindow should return INFO, WARN, and ERROR totals for each time bucket
// across all requests.
func LevelsByWindow(records []LogRecord, bucketSize BucketSize) []LevelVolumePoint {
	panic("not implemented")
}

// WarnAndErrorCountsByWindow should return WARN and ERROR totals for each time bucket.
// INFO counts may be zeroed or ignored depending on the final implementation.
func WarnAndErrorCountsByWindow(records []LogRecord, bucketSize BucketSize) []LevelVolumePoint {
	panic("not implemented")
}

// StatusClassesByWindow should return totals for each HTTP status class
// (1xx, 2xx, 3xx, 4xx, 5xx) for every time bucket.
func StatusClassesByWindow(records []LogRecord, bucketSize BucketSize) []StatusClassVolumePoint {
	panic("not implemented")
}

// StatusCodesByWindow should return totals for each exact HTTP status code
// seen in each time bucket.
// This is intended to support views such as:
// 200=152, 201=18, 404=6, 429=2, 500=1.
func StatusCodesByWindow(records []LogRecord, bucketSize BucketSize) []StatusCodeVolumePoint {
	panic("not implemented")
}

// MetricsByPath should return aggregate metrics for each distinct path across the
// full input set, regardless of time bucket.
// This is intended to support questions like:
// "How many requests hit /api/orders?"
// "How many WARN/ERROR logs came from /api/login?"
// "What was the average/max latency for /api/reports/daily?"
//
// It is expected to delegate its core aggregation work to aggregatePathMetrics.
func MetricsByPath(records []LogRecord) []PathMetrics {
	// panic("not implemented")
	returnMetrics := aggregatePathMetrics(records)

	slices.SortFunc(returnMetrics, func(a, b PathMetrics) int {
		return strings.Compare(a.Path, b.Path)
	})

	return returnMetrics
}

// aggregatePathMetrics should contain the shared "group by path and aggregate"
// behavior used by both MetricsByPath and MetricsByPathAndWindow.
// It should return one PathMetrics value per distinct path in the provided records.
func aggregatePathMetrics(records []LogRecord) []PathMetrics {
	// panic("not implemented")

	var metrics []PathMetrics
	var seenPaths = make(map[string]PathMetrics, len(records))

	for i := range records {
		record := records[i]

		_, ok := seenPaths[record.Path]

		if !ok {
			var tempLevelCount LevelCounts
			var tempStatusClassCounts StatusClassCounts
			var tempLatencySummary = LatencySummary{
				Count:   1,
				TotalMs: record.DurationMS,
				MaxMS:   record.DurationMS,
			}

			// Latency
			if record.DurationMS >= 500 {
				tempLatencySummary.SlowOver500MS += 1
				tempLatencySummary.SlowOver250MS += 1
				tempLatencySummary.SlowOver100MS += 1
			} else if record.DurationMS >= 250 {
				tempLatencySummary.SlowOver250MS += 1
				tempLatencySummary.SlowOver100MS += 1
			} else if record.DurationMS >= 100 {
				tempLatencySummary.SlowOver100MS += 1
			}

			// Level
			switch record.Level {
			case "INFO":
				tempLevelCount.InfoCount = 1
			case "WARN":
				tempLevelCount.WarnCount = 1
			case "ERROR":
				tempLevelCount.ErrorCount = 1
			}

			// Status
			switch {
			case record.Status < 200:
				tempStatusClassCounts.Status1xx = 1
			case record.Status < 300:
				tempStatusClassCounts.Status2xx = 1
			case record.Status < 400:
				tempStatusClassCounts.Status3xx = 1
			case record.Status < 500:
				tempStatusClassCounts.Status4xx = 1
			case record.Status < 600:
				tempStatusClassCounts.Status5xx = 1
			}

			seenPaths[record.Path] = PathMetrics{
				Path:         record.Path,
				RequestCount: 1,
				LevelCounts:  tempLevelCount,
				StatusCounts: tempStatusClassCounts,
				Latency:      tempLatencySummary,
			}
		} else {
			currentData := seenPaths[record.Path]
			currentData.RequestCount += 1

			// Latency
			if record.DurationMS >= 500 {
				currentData.Latency.SlowOver500MS += 1
				currentData.Latency.SlowOver250MS += 1
				currentData.Latency.SlowOver100MS += 1
			} else if record.DurationMS >= 250 {
				currentData.Latency.SlowOver250MS += 1
				currentData.Latency.SlowOver100MS += 1
			} else if record.DurationMS >= 100 {
				currentData.Latency.SlowOver100MS += 1
			}

			if record.DurationMS > currentData.Latency.MaxMS {
				currentData.Latency.MaxMS = record.DurationMS
			}

			currentData.Latency.Count += 1
			currentData.Latency.TotalMs += record.DurationMS

			// Level
			switch record.Level {
			case "INFO":
				currentData.LevelCounts.InfoCount += 1
			case "WARN":
				currentData.LevelCounts.WarnCount += 1
			case "ERROR":
				currentData.LevelCounts.ErrorCount += 1
			}

			// Status
			switch {
			case record.Status < 200:
				currentData.StatusCounts.Status1xx += 1
			case record.Status < 300:
				currentData.StatusCounts.Status2xx += 1
			case record.Status < 400:
				currentData.StatusCounts.Status3xx += 1
			case record.Status < 500:
				currentData.StatusCounts.Status4xx += 1
			case record.Status < 600:
				currentData.StatusCounts.Status5xx += 1
			}
			seenPaths[record.Path] = currentData
		}
	}

	for _, value := range seenPaths {
		value.Latency.AverageMS = value.Latency.TotalMs / value.Latency.Count
		metrics = append(metrics, value)
	}

	return metrics
}

// groupRecordsByWindow should contain the shared "split records into contiguous
// time buckets" behavior used by windowed metric functions.
// It should return one PathWindowMetrics value per bucket, with each bucket
// containing the records that fall within that bucket's time range.
func groupRecordsByWindow(records []LogRecord, bucketSize BucketSize) []PathWindowMetrics {
	panic("not implemented")
}

// MetricsByPathAndWindow should return per-path metrics inside each time bucket.
// This is the most dashboard-friendly aggregation because it combines:
// time window + path + request count + level totals + status totals + latency.
func MetricsByPathAndWindow(records []LogRecord, bucketSize BucketSize) []PathWindowMetrics {
	// panic("not implemented")

	if len(records) == 0 {
		return []PathWindowMetrics{}
	}

	slices.SortFunc(records, func(a, b LogRecord) int {
		if a.TS.Unix() > b.TS.Unix() {
			return 1
		} else if a.TS.Unix() < b.TS.Unix() {
			return -1
		} else {
			return 0
		}
	})

	endTime := records[0].TS.Add(time.Duration(bucketSize))
	var returnPathWindowMetrics []PathWindowMetrics
	var tempLogRecords []LogRecord

	// Loop through log records
	for _, record := range records {

		if len(tempLogRecords) == 0 {
			tempLogRecords = append(tempLogRecords, record)
		} else if record.TS.Unix() <= endTime.Unix() {
			tempLogRecords = append(tempLogRecords, record)
		} else {

			if len(tempLogRecords) > 0 {
				returnPathWindowMetrics = append(returnPathWindowMetrics, PathWindowMetrics{
					Window: TimeWindow{
						Start: tempLogRecords[0].TS,
						End:   tempLogRecords[len(tempLogRecords)-1].TS,
					},
					Paths: aggregatePathMetrics(tempLogRecords),
				})
			}

			endTime = record.TS.Add(time.Duration(bucketSize))
			tempLogRecords = []LogRecord{record}
		}
	}

	if len(tempLogRecords) > 0 {
		returnPathWindowMetrics = append(returnPathWindowMetrics, PathWindowMetrics{
			Window: TimeWindow{
				Start: tempLogRecords[0].TS,
				End:   tempLogRecords[len(tempLogRecords)-1].TS,
			},
			Paths: aggregatePathMetrics(tempLogRecords),
		})
	}

	return returnPathWindowMetrics
}

// LatencyByPath should return latency summaries for each path across the full input set.
// It is intended to answer:
// "Which API paths are slow on average?"
// "Which API path had the highest max latency?"
func LatencyByPath(records []LogRecord) []PathLatencyMetrics {
	panic("not implemented")
}

// SlowRequestsByPath should return per-path counts of requests that exceed configured
// thresholds such as 100ms, 250ms, and 500ms.
func SlowRequestsByPath(records []LogRecord) []PathLatencyMetrics {
	panic("not implemented")
}

// SlowRequestsByWindow should return the number of slow requests inside each time bucket.
// Depending on the final design, this may represent:
// all requests over a chosen threshold, or
// a combined count of threshold breaches.
func SlowRequestsByWindow(records []LogRecord, bucketSize BucketSize) []RequestVolumePoint {
	panic("not implemented")
}

// ErrorRateByPath should return per-path request totals alongside status-class counts
// so error rates can be calculated for each route.
// This is intended to support views like:
// "/api/checkout had 120 requests, 6 of which were 5xx."
func ErrorRateByPath(records []LogRecord) []PathMetrics {
	panic("not implemented")
}

// ErrorRateByWindow should return status-class totals per time bucket so overall
// error rates can be calculated over time.
// This is useful for identifying windows with elevated 4xx or 5xx traffic.
func ErrorRateByWindow(records []LogRecord, bucketSize BucketSize) []StatusClassVolumePoint {
	panic("not implemented")
}
