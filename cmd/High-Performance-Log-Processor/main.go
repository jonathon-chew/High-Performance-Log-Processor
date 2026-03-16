package main

import (
	"High-Performance-Log-Processor/internal/cli"
	"High-Performance-Log-Processor/internal/dashboard"
	"High-Performance-Log-Processor/internal/parseinput"
	"High-Performance-Log-Processor/internal/ping"
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// ts=2026-03-14T09:01:20.006Z level=INFO req_id=0f4c9f3d method=GET path=/health status=200 duration_ms=1 bytes=2 ip=10.0.0.5 ua="kube-probe/1.31" msg="request complete"

func main() {

	var Flags cli.Flags

	if len(os.Args) > 1 {
		Flags = cli.CLI(os.Args[1:])
	}

	if Flags.Ping {
		ping.ParsePing()
	}

	if Flags.FileName == "" {
		os.Exit(1)
	}

	filePointer, err := os.Open(Flags.FileName)
	if err != nil {
		log.Print("[ERROR]: Unable to open file ", err)
		return
	}

	defer filePointer.Close()

	bufReader := bufio.NewScanner(filePointer)

	for bufReader.Scan() {
		line := bufReader.Text()

		var splitLog []string
		var insideQuote bool
		var start int
		for curByte := range line {
			if line[curByte] == '"' {
				if insideQuote {
					insideQuote = false
				} else {
					insideQuote = true
				}
			}

			if line[curByte] == ' ' && !insideQuote {
				splitLog = append(splitLog, line[start:curByte])
				start = curByte + 1
			}
		}

		// add the final message!
		splitLog = append(splitLog, line[start:])

		parseinput.Logs = append(parseinput.Logs, parseinput.BuildLogRecord(splitLog))
	}

	if bufReader.Err() != nil {
		log.Panic(bufReader.Err())
		return
	}

	for _, i := range dashboard.MetricsByPath(parseinput.Logs) {
		err := json.NewEncoder(os.Stdout).Encode(i)
		if err != nil {
			continue
		}
	}
}
