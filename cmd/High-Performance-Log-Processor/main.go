package main

import (
	"High-Performance-Log-Processor/internal/cli"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ts=2026-03-14T09:01:20.006Z level=INFO req_id=0f4c9f3d method=GET path=/health status=200 duration_ms=1 bytes=2 ip=10.0.0.5 ua="kube-probe/1.31" msg="request complete"
type LogLine struct {
	TS        string
	Level     string
	RequestId string
	Method    string
	Path      string
	Status    string
	Duration  string
	Bytes     string
	Ip        string
	UserAgent string
	Message   string
}

type Logs struct {
	AllLogs []LogLine
}

func main() {

	var Flags cli.Flags
	var Logs Logs

	if len(os.Args) > 1 {
		Flags = cli.CLI(os.Args[1:])
	}

	if Flags.FileName == "" {
		os.Exit(1)
	}

	filePointer, err := os.Open(Flags.FileName)
	if err != nil {
		log.Print("[ERROR]: Unable to open file ", err)
	}

	defer filePointer.Close()

	bufReader := bufio.NewScanner(filePointer)

	for bufReader.Scan() {
		line := bufReader.Text()
		// fmt.Print(line, "*\n*")

		var splitLog []string
		var insideQuote bool
		var start int
		// fmt.Print("Split Log Length: ", len(splitLog), "\n")
		for curByte := range line {
			if line[curByte] == '"' {
				if insideQuote {
					insideQuote = false
				} else {
					insideQuote = true
					log.Print("flip")
				}
			}

			if line[curByte] == ' ' && !insideQuote {
				splitLog = append(splitLog, line[start:curByte])
				start = curByte + 1
			}
		}

		// add the final message!
		splitLog = append(splitLog, line[start:])

		fmt.Print("Split log length: ", len(splitLog), "\n")
		for i := range splitLog {
			fmt.Print(i, " ", splitLog[i], "\n")
		}

		Logs.AllLogs = append(Logs.AllLogs, LogLine{
			TS:        strings.Split(splitLog[0], "=")[1],
			Level:     strings.Split(splitLog[1], "=")[1],
			RequestId: strings.Split(splitLog[2], "=")[1],
			Method:    strings.Split(splitLog[3], "=")[1],
			Path:      strings.Split(splitLog[4], "=")[1],
			Status:    strings.Split(splitLog[5], "=")[1],
			Duration:  strings.Split(splitLog[6], "=")[1],
			Bytes:     strings.Split(splitLog[7], "=")[1],
			Ip:        strings.Split(splitLog[8], "=")[1],
			UserAgent: strings.Split(splitLog[9], "=")[1],
			Message:   strings.Split(splitLog[10], "=")[1],
		})
	}

	var warnings int

	for _, log := range Logs.AllLogs {
		if log.Level == "WARN" {
			warnings += 1
		}
	}
	fmt.Print("There were: ", len(Logs.AllLogs), " logs. ", strconv.Itoa(warnings), " were warnings")
}
