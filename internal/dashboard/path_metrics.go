package dashboard

import (
	"slices"
	"strings"
)

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

// LatencyByPath should return latency summaries for each path across the full input set.
// It is intended to answer:
// "Which API paths are slow on average?"
// "Which API path had the highest max latency?"
func LatencyByPath(records []LogRecord) []PathLatencyMetrics {
	// panic("not implemented")

	var returnPathLatencyMetrics []PathLatencyMetrics

	data := aggregatePathMetrics(records)

	for _, record := range data {
		returnPathLatencyMetrics = append(returnPathLatencyMetrics, PathLatencyMetrics{
			Path:    record.Path,
			Latency: record.Latency,
		})
	}

	slices.SortFunc(returnPathLatencyMetrics, func(a, b PathLatencyMetrics) int {
		return strings.Compare(a.Path, b.Path)
	})

	return returnPathLatencyMetrics
}

// SlowRequestsByPath should return per-path counts of requests that exceed configured
// thresholds such as 100ms, 250ms, and 500ms.
func SlowRequestsByPath(records []LogRecord) []PathLatencyMetrics {
	// panic("not implemented")

	var returnPathLatencyMetrics []PathLatencyMetrics

	if len(records) == 0 {
		return []PathLatencyMetrics{}
	}

	data := aggregatePathMetrics(records)

	for _, eachPath := range data {
		if eachPath.Latency.SlowOver100MS > 0 || eachPath.Latency.SlowOver250MS > 0 || eachPath.Latency.SlowOver500MS > 0 {
			returnPathLatencyMetrics = append(returnPathLatencyMetrics, PathLatencyMetrics{
				Path:    eachPath.Path,
				Latency: eachPath.Latency,
			})
		}
	}

	slices.SortFunc(returnPathLatencyMetrics, func(a, b PathLatencyMetrics) int {
		return strings.Compare(a.Path, b.Path)
	})

	return returnPathLatencyMetrics
}

// ErrorRateByPath should return per-path request totals alongside status-class counts
// so error rates can be calculated for each route.
// This is intended to support views like:
// "/api/checkout had 120 requests, 6 of which were 5xx."
func ErrorRateByPath(records []LogRecord) []PathMetrics {
	// panic("not implemented")

	if len(records) == 0 {
		return []PathMetrics{}
	}

	data := aggregatePathMetrics(records)

	slices.SortFunc(data, func(a, b PathMetrics) int {
		return strings.Compare(a.Path, b.Path)
	})

	return data
}
