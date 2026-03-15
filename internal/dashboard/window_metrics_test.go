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
