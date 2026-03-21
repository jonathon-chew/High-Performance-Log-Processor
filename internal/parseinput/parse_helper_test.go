package parseinput

import (
	"High-Performance-Log-Processor/internal/cli"
	"testing"
	"time"
)

func TestBuildLogRecordKeepsKnownFieldsWhenMalformedTokensExist(t *testing.T) {
	record := BuildLogRecord([]string{
		`ts=not-a-time`,
		`level=WARN`,
		`path=/api/test`,
		`status=nope`,
		`duration_ms=bad`,
		`bytes=broken`,
		`ua="curl/8.7.1"`,
		`msg="request complete"`,
		`broken_token`,
	})

	if !record.TS.IsZero() {
		t.Fatalf("expected invalid timestamp to leave zero time, got %v", record.TS)
	}
	if record.Level != "WARN" {
		t.Fatalf("expected level WARN, got %q", record.Level)
	}
	if record.Path != "/api/test" {
		t.Fatalf("expected path /api/test, got %q", record.Path)
	}
	if record.Status != 0 || record.DurationMS != 0 || record.Bytes != 0 {
		t.Fatalf("expected invalid numeric fields to fall back to zero, got %+v", record)
	}
	if record.UserAgent != "curl/8.7.1" {
		t.Fatalf("expected user agent to be unquoted, got %q", record.UserAgent)
	}
	if record.Message != "request complete" {
		t.Fatalf("expected message to be unquoted, got %q", record.Message)
	}
}

func TestParseFileLargeFixtureReturnsExpectedRecordCount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large fixture test in short mode")
	}

	Logs = nil
	t.Cleanup(func() {
		Logs = nil
	})

	records := ParseFile(cli.Flags{FileName: "../../testdata/access-large.log"})
	if len(records) != 1000000 {
		t.Fatalf("expected 1000000 records from large fixture, got %d", len(records))
	}
	if records[0].Path == "" {
		t.Fatal("expected first parsed record to contain a path")
	}
	if records[len(records)-1].Path == "" {
		t.Fatal("expected last parsed record to contain a path")
	}
}

func TestBuildLogRecordPreservesValidTimestampAlongsideMalformedFields(t *testing.T) {
	record := BuildLogRecord([]string{
		`ts=2026-03-14T09:01:20.006Z`,
		`status=broken`,
		`duration_ms=oops`,
		`path=/api/orders`,
	})

	if record.TS.IsZero() {
		t.Fatal("expected valid timestamp to parse")
	}
	if !record.TS.Equal(time.Date(2026, 3, 14, 9, 1, 20, 6*1000000, time.UTC)) {
		t.Fatalf("unexpected parsed time: %v", record.TS)
	}
	if record.Path != "/api/orders" {
		t.Fatalf("expected path /api/orders, got %q", record.Path)
	}
}

