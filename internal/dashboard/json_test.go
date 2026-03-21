package dashboard

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPathMetricsJSONUsesTaggedFieldNames(t *testing.T) {
	data, err := json.Marshal(PathMetrics{
		Path:         "/api/login",
		RequestCount: 2,
		LevelCounts: LevelCounts{
			WarnCount: 1,
		},
		StatusCounts: StatusClassCounts{
			Status4xx: 1,
		},
		Latency: LatencySummary{
			Count:         2,
			TotalMs:       450,
			AverageMS:     225,
			MaxMS:         300,
			SlowOver100MS: 2,
			SlowOver250MS: 1,
		},
	})
	if err != nil {
		t.Fatalf("marshal path metrics: %v", err)
	}

	jsonText := string(data)
	expectedFragments := []string{
		`"path":"/api/login"`,
		`"request_count":2`,
		`"level_counts":{"warn_count":1}`,
		`"status_counts":{"status_4xx":1}`,
		`"latency":{"count":2,"total_ms":450,"average_ms":225,"max_ms":300,"slow_over_100_ms":2,"slow_over_250_ms":1}`,
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(jsonText, fragment) {
			t.Fatalf("expected JSON to contain %s, got %s", fragment, jsonText)
		}
	}
}

func TestStatusCodeVolumePointJSONUsesTaggedFieldNames(t *testing.T) {
	data, err := json.Marshal(StatusCodeVolumePoint{
		Window: TimeWindow{
			Start: mustTestTime("2026-03-14T09:00:00Z"),
			End:   mustTestTime("2026-03-14T09:05:00Z"),
		},
		Counts: []StatusCodeCount{
			{StatusCode: 200, Count: 10},
			{StatusCode: 404, Count: 2},
		},
	})
	if err != nil {
		t.Fatalf("marshal status code volume point: %v", err)
	}

	jsonText := string(data)
	expectedFragments := []string{
		`"window":{"start":"2026-03-14T09:00:00Z","end":"2026-03-14T09:05:00Z"}`,
		`"counts":[{"status_code":200,"count":10},{"status_code":404,"count":2}]`,
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(jsonText, fragment) {
			t.Fatalf("expected JSON to contain %s, got %s", fragment, jsonText)
		}
	}
}

