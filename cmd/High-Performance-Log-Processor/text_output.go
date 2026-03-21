package main

import (
	"fmt"
	"strings"

	"High-Performance-Log-Processor/internal/dashboard"
)

const pathColumnWidth = 24

func printTextOutput(record any) {
	switch value := record.(type) {
	case dashboard.PathMetrics:
		fmt.Printf(
			"%-24s req=%6d info=%4d warn=%4d err=%4d 2xx=%4d 4xx=%4d 5xx=%4d avg_ms=%4d max_ms=%4d\n",
			value.Path,
			value.RequestCount,
			value.LevelCounts.InfoCount,
			value.LevelCounts.WarnCount,
			value.LevelCounts.ErrorCount,
			value.StatusCounts.Status2xx,
			value.StatusCounts.Status4xx,
			value.StatusCounts.Status5xx,
			value.Latency.AverageMS,
			value.Latency.MaxMS,
		)
	case dashboard.PathLatencyMetrics:
		fmt.Printf(
			"%-24s avg_ms=%4d max_ms=%4d over_100=%4d over_250=%4d over_500=%4d\n",
			value.Path,
			value.Latency.AverageMS,
			value.Latency.MaxMS,
			value.Latency.SlowOver100MS,
			value.Latency.SlowOver250MS,
			value.Latency.SlowOver500MS,
		)
	case dashboard.RequestVolumePoint:
		fmt.Printf("%-43s requests=%6d\n", formatWindow(value.Window), value.RequestCount)
	case dashboard.LevelVolumePoint:
		fmt.Printf(
			"%-43s info=%4d warn=%4d err=%4d\n",
			formatWindow(value.Window),
			value.Counts.InfoCount,
			value.Counts.WarnCount,
			value.Counts.ErrorCount,
		)
	case dashboard.StatusClassVolumePoint:
		fmt.Printf(
			"%-43s 1xx=%4d 2xx=%4d 3xx=%4d 4xx=%4d 5xx=%4d\n",
			formatWindow(value.Window),
			value.Counts.Status1xx,
			value.Counts.Status2xx,
			value.Counts.Status3xx,
			value.Counts.Status4xx,
			value.Counts.Status5xx,
		)
	case dashboard.StatusCodeVolumePoint:
		fmt.Printf("%-43s %s\n", formatWindow(value.Window), formatStatusCodeCounts(value.Counts))
	case dashboard.PathWindowMetrics:
		fmt.Printf(
			"%-43s paths=%4d requests=%6d\n",
			formatWindow(value.Window),
			len(value.Paths),
			totalRequests(value.Paths),
		)
	default:
		fmt.Println(record)
	}
}

func formatWindow(window dashboard.TimeWindow) string {
	return fmt.Sprintf("%s -> %s", window.Start.Format("2006-01-02T15:04:05Z07:00"), window.End.Format("2006-01-02T15:04:05Z07:00"))
}

func formatStatusCodeCounts(counts []dashboard.StatusCodeCount) string {
	if len(counts) == 0 {
		return "no_status_codes"
	}

	var parts []string
	for _, count := range counts {
		parts = append(parts, fmt.Sprintf("%d=%d", count.StatusCode, count.Count))
	}
	return strings.Join(parts, " ")
}

func totalRequests(paths []dashboard.PathMetrics) int {
	var total int
	for _, path := range paths {
		total += path.RequestCount
	}
	return total
}

