package dashboard

import "testing"

func TestMetricsByPathAggregatesByPath(t *testing.T) {
	metrics := MetricsByPath(sampleRecords())
	if len(metrics) != 3 {
		t.Fatalf("expected 3 path metrics, got %d", len(metrics))
	}

	login := metrics[0]
	if login.Path != "/api/login" {
		t.Fatalf("expected first path to be /api/login, got %q", login.Path)
	}
	if login.RequestCount != 2 {
		t.Fatalf("expected login request count 2, got %d", login.RequestCount)
	}
	if login.LevelCounts.InfoCount != 1 || login.LevelCounts.WarnCount != 1 || login.LevelCounts.ErrorCount != 0 {
		t.Fatalf("unexpected login level counts: %+v", login.LevelCounts)
	}
	if login.StatusCounts.Status2xx != 1 || login.StatusCounts.Status4xx != 1 {
		t.Fatalf("unexpected login status counts: %+v", login.StatusCounts)
	}
	if login.Latency.AverageMS != 225 || login.Latency.MaxMS != 300 {
		t.Fatalf("unexpected login latency summary: %+v", login.Latency)
	}
	if login.Latency.SlowOver100MS != 2 || login.Latency.SlowOver250MS != 1 || login.Latency.SlowOver500MS != 0 {
		t.Fatalf("unexpected login slow counts: %+v", login.Latency)
	}

	orders := metrics[1]
	if orders.Path != "/api/orders" {
		t.Fatalf("expected second path to be /api/orders, got %q", orders.Path)
	}
	if orders.RequestCount != 2 || orders.StatusCounts.Status2xx != 2 {
		t.Fatalf("unexpected orders metrics: %+v", orders)
	}
	if orders.Latency.AverageMS != 35 || orders.Latency.MaxMS != 50 {
		t.Fatalf("unexpected orders latency summary: %+v", orders.Latency)
	}

	reports := metrics[2]
	if reports.Path != "/api/reports" {
		t.Fatalf("expected third path to be /api/reports, got %q", reports.Path)
	}
	if reports.LevelCounts.ErrorCount != 1 || reports.StatusCounts.Status5xx != 1 {
		t.Fatalf("unexpected reports metrics: %+v", reports)
	}
	if reports.Latency.AverageMS != 600 || reports.Latency.MaxMS != 600 {
		t.Fatalf("unexpected reports latency summary: %+v", reports.Latency)
	}
	if reports.Latency.SlowOver100MS != 1 || reports.Latency.SlowOver250MS != 1 || reports.Latency.SlowOver500MS != 1 {
		t.Fatalf("unexpected reports slow counts: %+v", reports.Latency)
	}
}

func TestLatencyByPathProjectsLatencyOnly(t *testing.T) {
	metrics := LatencyByPath(sampleRecords())
	if len(metrics) != 3 {
		t.Fatalf("expected 3 latency metrics, got %d", len(metrics))
	}
	if metrics[0].Path != "/api/login" || metrics[1].Path != "/api/orders" || metrics[2].Path != "/api/reports" {
		t.Fatalf("unexpected sorted paths: %+v", metrics)
	}
	if metrics[2].Latency.MaxMS != 600 {
		t.Fatalf("expected reports max latency 600, got %+v", metrics[2].Latency)
	}
}

func TestSlowRequestsByPathReturnsOnlySlowPaths(t *testing.T) {
	metrics := SlowRequestsByPath(sampleRecords())
	if len(metrics) != 2 {
		t.Fatalf("expected 2 slow paths, got %d", len(metrics))
	}
	if metrics[0].Path != "/api/login" || metrics[1].Path != "/api/reports" {
		t.Fatalf("unexpected slow paths: %+v", metrics)
	}
}

func TestErrorRateByPathReusesPathAggregation(t *testing.T) {
	metrics := ErrorRateByPath(sampleRecords())
	if len(metrics) != 3 {
		t.Fatalf("expected 3 path metrics, got %d", len(metrics))
	}
	if metrics[0].Path != "/api/login" || metrics[2].Path != "/api/reports" {
		t.Fatalf("unexpected sorted paths: %+v", metrics)
	}
	if metrics[2].StatusCounts.Status5xx != 1 {
		t.Fatalf("expected reports to have one 5xx, got %+v", metrics[2].StatusCounts)
	}
}
