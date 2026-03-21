package cli

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func captureOutput(t *testing.T, fn func()) (string, string) {
	t.Helper()

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stderr pipe: %v", err)
	}

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	fn()

	_ = stdoutWriter.Close()
	_ = stderrWriter.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdoutBytes, err := io.ReadAll(stdoutReader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	stderrBytes, err := io.ReadAll(stderrReader)
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}

	return string(stdoutBytes), string(stderrBytes)
}

func TestCLIReturnsFileNameForDefaultFileInput(t *testing.T) {
	flags := CLI([]string{"../../testdata/access.log"})

	if flags.FileName != "../../testdata/access.log" {
		t.Fatalf("expected file name to be preserved, got %q", flags.FileName)
	}
	if flags.Ping {
		t.Fatal("did not expect ping mode for file input")
	}
}

func TestCLIParsesOutputFlagAsJSON(t *testing.T) {
	flags := CLI([]string{"../../testdata/access.log", "--output", "json"})

	if flags.Output != "JSON" {
		t.Fatalf("expected output JSON, got %q", flags.Output)
	}
}

func TestCLIParsesBucketDuration(t *testing.T) {
	flags := CLI([]string{"../../testdata/access.log", "RequestsByWindow", "--time", "10m"})

	if flags.Bucket != 10*time.Minute {
		t.Fatalf("expected 10m bucket, got %v", flags.Bucket)
	}
	if !flags.RequestsByWindow {
		t.Fatal("expected RequestsByWindow flag to be set")
	}
}

func TestCLIMissingTimeValueDefaultsToFiveMinutes(t *testing.T) {
	flags := CLI([]string{"../../testdata/access.log", "RequestsByWindow", "--time"})

	if flags.Bucket != 5*time.Minute {
		t.Fatalf("expected default 5m bucket, got %v", flags.Bucket)
	}
}

func TestCLIMissingOutputValueLeavesOutputUnset(t *testing.T) {
	flags := CLI([]string{"../../testdata/access.log", "--output"})

	if flags.Output != "" {
		t.Fatalf("expected output to remain unset, got %q", flags.Output)
	}
}

func TestCLIParsesSpecificMetricSelection(t *testing.T) {
	flags := CLI([]string{"../../testdata/access.log", "ErrorRateByWindow"})

	if !flags.ErrorRateByWindow {
		t.Fatal("expected ErrorRateByWindow flag to be set")
	}
	if flags.MetricsByPath {
		t.Fatal("did not expect MetricsByPath to be set")
	}
}

func TestCLIReturnsPingImmediately(t *testing.T) {
	flags := CLI([]string{"ping", "MetricsByPath", "--output", "JSON"})

	if !flags.Ping {
		t.Fatal("expected ping mode to be enabled")
	}
	if flags.MetricsByPath {
		t.Fatal("did not expect later args to be parsed after ping")
	}
	if flags.Output != "" {
		t.Fatalf("expected no output format after ping early return, got %q", flags.Output)
	}
}

func TestCLIHelpPrintsUsageToStdout(t *testing.T) {
	stdout, stderr := captureOutput(t, func() {
		flags := CLI([]string{"help"})
		if flags != (Flags{}) {
			t.Fatalf("expected empty flags for help, got %+v", flags)
		}
	})

	if !strings.Contains(stdout, "USAGE:") {
		t.Fatalf("expected usage text on stdout, got %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("expected no stderr for help, got %q", stderr)
	}
}

func TestCLIVersionPrintsVersionToStdout(t *testing.T) {
	stdout, stderr := captureOutput(t, func() {
		flags := CLI([]string{"version"})
		if flags != (Flags{}) {
			t.Fatalf("expected empty flags for version, got %+v", flags)
		}
	})

	if !strings.Contains(stdout, version) {
		t.Fatalf("expected version output to contain %q, got %q", version, stdout)
	}
	if stderr != "" {
		t.Fatalf("expected no stderr for version, got %q", stderr)
	}
}
