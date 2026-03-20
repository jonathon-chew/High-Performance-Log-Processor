package dashboard

import "testing"

func TestGroupRecordsByWindowBucketsSortedRecords(t *testing.T) {
	buckets := groupRecordsByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}

	if len(buckets[0].Records) != 3 {
		t.Fatalf("expected first bucket to have 3 records, got %d", len(buckets[0].Records))
	}
	if got := buckets[0].Window.Start; !got.Equal(mustTestTime("2026-03-14T09:00:00Z")) {
		t.Fatalf("unexpected first bucket start: %v", got)
	}
	if got := buckets[0].Window.End; !got.Equal(mustTestTime("2026-03-14T09:04:00Z")) {
		t.Fatalf("unexpected first bucket end: %v", got)
	}

	if len(buckets[1].Records) != 2 {
		t.Fatalf("expected second bucket to have 2 records, got %d", len(buckets[1].Records))
	}
	if got := buckets[1].Window.Start; !got.Equal(mustTestTime("2026-03-14T09:06:00Z")) {
		t.Fatalf("unexpected second bucket start: %v", got)
	}
	if got := buckets[1].Window.End; !got.Equal(mustTestTime("2026-03-14T09:07:00Z")) {
		t.Fatalf("unexpected second bucket end: %v", got)
	}
}

func TestMetricsByPathAndWindowAggregatesEachBucket(t *testing.T) {
	metrics := MetricsByPathAndWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 window metrics, got %d", len(metrics))
	}

	if len(metrics[0].Paths) != 2 {
		t.Fatalf("expected 2 path metrics in first window, got %d", len(metrics[0].Paths))
	}
	if len(metrics[1].Paths) != 2 {
		t.Fatalf("expected 2 path metrics in second window, got %d", len(metrics[1].Paths))
	}
}

func TestRequestsByWindowCountsBucketLengths(t *testing.T) {
	metrics := RequestsByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 request-volume windows, got %d", len(metrics))
	}

	if metrics[0].RequestCount != 3 {
		t.Fatalf("expected first window request count 3, got %d", metrics[0].RequestCount)
	}
	if metrics[1].RequestCount != 2 {
		t.Fatalf("expected second window request count 2, got %d", metrics[1].RequestCount)
	}
}

func TestLevelsByWindowSumsAllLevelsPerBucket(t *testing.T) {
	metrics := LevelsByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 level windows, got %d", len(metrics))
	}

	first := metrics[0].Counts
	if first.InfoCount != 2 || first.WarnCount != 1 || first.ErrorCount != 0 {
		t.Fatalf("unexpected first window level totals: %+v", first)
	}

	second := metrics[1].Counts
	if second.InfoCount != 1 || second.WarnCount != 0 || second.ErrorCount != 1 {
		t.Fatalf("unexpected second window level totals: %+v", second)
	}
}

func TestWarnAndErrorCountsByWindowReturnsPerWindowTotals(t *testing.T) {
	metrics := WarnAndErrorCountsByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 warn/error windows, got %d", len(metrics))
	}

	first := metrics[0].Counts
	if first.WarnCount != 1 || first.ErrorCount != 0 {
		t.Fatalf("unexpected first window warn/error totals: %+v", first)
	}

	second := metrics[1].Counts
	if second.WarnCount != 0 || second.ErrorCount != 1 {
		t.Fatalf("unexpected second window warn/error totals: %+v", second)
	}
}

func TestStatusClassesByWindowSumsStatusClassesPerBucket(t *testing.T) {
	metrics := StatusClassesByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 status-class windows, got %d", len(metrics))
	}

	first := metrics[0].Counts
	if first.Status2xx != 2 || first.Status4xx != 1 || first.Status5xx != 0 {
		t.Fatalf("unexpected first window status-class totals: %+v", first)
	}

	second := metrics[1].Counts
	if second.Status2xx != 1 || second.Status4xx != 0 || second.Status5xx != 1 {
		t.Fatalf("unexpected second window status-class totals: %+v", second)
	}
}

func TestStatusCodesByWindowCountsExactCodesPerBucket(t *testing.T) {
	metrics := StatusCodesByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 status-code windows, got %d", len(metrics))
	}

	first := metrics[0].Counts
	if len(first) != 3 {
		t.Fatalf("expected 3 status codes in first window, got %d", len(first))
	}
	if first[0].StatusCode != 200 || first[0].Count != 2 {
		t.Fatalf("unexpected first status code entry: %+v", first[0])
	}
	if first[1].StatusCode != 401 || first[1].Count != 1 {
		t.Fatalf("unexpected second status code entry: %+v", first[1])
	}

	second := metrics[1].Counts
	if len(second) != 2 {
		t.Fatalf("expected 2 status codes in second window, got %d", len(second))
	}
	if second[0].StatusCode != 201 || second[0].Count != 1 {
		t.Fatalf("unexpected first status code in second window: %+v", second[0])
	}
	if second[1].StatusCode != 503 || second[1].Count != 1 {
		t.Fatalf("unexpected second status code in second window: %+v", second[1])
	}
}

func TestSlowRequestsByWindowCountsRequestsAboveThreshold(t *testing.T) {
	metrics := SlowRequestsByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 slow-request windows, got %d", len(metrics))
	}

	if metrics[0].RequestCount != 2 {
		t.Fatalf("expected first window slow-request count 2, got %d", metrics[0].RequestCount)
	}
	if metrics[1].RequestCount != 1 {
		t.Fatalf("expected second window slow-request count 1, got %d", metrics[1].RequestCount)
	}
}

func TestErrorRateByWindowSumsStatusCountsAcrossPaths(t *testing.T) {
	metrics := ErrorRateByWindow(sampleRecords(), BucketSize(5*60*1000000000))
	if len(metrics) != 2 {
		t.Fatalf("expected 2 error-rate windows, got %d", len(metrics))
	}

	first := metrics[0].Counts
	if first.Status2xx != 2 || first.Status4xx != 1 || first.Status5xx != 0 {
		t.Fatalf("unexpected first window status totals: %+v", first)
	}

	second := metrics[1].Counts
	if second.Status2xx != 1 || second.Status4xx != 0 || second.Status5xx != 1 {
		t.Fatalf("unexpected second window status totals: %+v", second)
	}
}
