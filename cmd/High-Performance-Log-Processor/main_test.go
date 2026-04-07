package main

import (
	"testing"

	"github.com/jonathon-chew/High-Performance-Log-Processor/internal/parseinput"
)

func TestStringToInt(t *testing.T) {
	if got := parseinput.StringToInt("123"); got != 123 {
		t.Fatalf("expected 123, got %d", got)
	}
	if got := parseinput.StringToInt("20.006Z"); got != 20 {
		t.Fatalf("expected 20 from decimal-like input, got %d", got)
	}
	if got := parseinput.StringToInt("nope"); got != 0 {
		t.Fatalf("expected 0 for invalid integer, got %d", got)
	}
}

func TestGetValueReturnsMatchingField(t *testing.T) {
	fields := []string{
		`level=INFO`,
		`path=/api/orders`,
		`ua="curl/8.7.1"`,
		`msg="request complete"`,
	}

	if got := parseinput.GetValue("path", fields); got != "/api/orders" {
		t.Fatalf("expected path value, got %q", got)
	}
	if got := parseinput.GetValue("ua", fields); got != "curl/8.7.1" {
		t.Fatalf("expected unquoted ua value, got %q", got)
	}
	if got := parseinput.GetValue("msg", fields); got != "request complete" {
		t.Fatalf("expected unquoted message value, got %q", got)
	}
}

func TestGetValueReturnsEmptyForMissingField(t *testing.T) {
	fields := []string{
		`level=INFO`,
		`broken_field`,
	}

	if got := parseinput.GetValue("msg", fields); got != "" {
		t.Fatalf("expected empty string for missing field, got %q", got)
	}
}

func TestParseTimeParsesTimestampParts(t *testing.T) {
	got := parseinput.ParseTime("ts=2026-03-14T09:01:20.006Z")
	if got.Year != 2026 || got.Month != 3 || got.Day != 14 {
		t.Fatalf("unexpected parsed date: %+v", got)
	}
	if got.Hour != 9 || got.Minute != 1 || got.Seconds != 20 {
		t.Fatalf("unexpected parsed time: %+v", got)
	}
}
