package main

import (
	"High-Performance-Log-Processor/internal/cli"
	"High-Performance-Log-Processor/internal/dashboard"
	"High-Performance-Log-Processor/internal/parseinput"
	"encoding/json"
	"os"
)

// ts=2026-03-14T09:01:20.006Z level=INFO req_id=0f4c9f3d method=GET path=/health status=200 duration_ms=1 bytes=2 ip=10.0.0.5 ua="kube-probe/1.31" msg="request complete"

func main() {

	var Flags cli.Flags
	var Logs []dashboard.LogRecord

	if len(os.Args) > 1 {
		Flags = cli.CLI(os.Args[1:])
	}

	if Flags.Ping {
		parseinput.ParsePing(Flags)
		return
	}

	if Flags.FileName == "" {
		os.Exit(1)
	} else {
		Logs = parseinput.ParseFile(Flags)
	}

	for _, i := range dashboard.MetricsByPath(Logs) {
		err := json.NewEncoder(os.Stdout).Encode(i)
		if err != nil {
			continue
		}
	}
}
