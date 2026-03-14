package main

import (
	"High-Performance-Log-Processor/internal/cli"
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

	Levels   []string
	Methods  []string
	Paths    []string
	Statuses []string
	IPs      []string
}

func GetValue(field string) string {
	splitLine := strings.Split(field, "=")
	// key := splitLine[0]
	value := splitLine[1]

	if value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-2]
	}

	return value
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

		if len(splitLog) != 11 {
			log.Panic("The line can not be parsed correctly")
		}

		var logLine = LogLine{
			TS:        GetValue(splitLog[0]),
			Level:     GetValue(splitLog[1]),
			RequestId: GetValue(splitLog[2]),
			Method:    GetValue(splitLog[3]),
			Path:      GetValue(splitLog[4]),
			Status:    GetValue(splitLog[5]),
			Duration:  GetValue(splitLog[6]),
			Bytes:     GetValue(splitLog[7]),
			Ip:        GetValue(splitLog[8]),
			// Hint: quoted fields like ua and msg still include their surrounding
			// quotes at this stage. Trimming those is a separate parsing step.
			UserAgent: GetValue(splitLog[9]),
			Message:   GetValue(splitLog[10]),
		}

		Logs.AllLogs = append(Logs.AllLogs, logLine)

		if !slices.Contains(Logs.Levels, logLine.Level) {
			Logs.Levels = append(Logs.Levels, logLine.Level)
		}
		if !slices.Contains(Logs.Methods, logLine.Method) {
			Logs.Methods = append(Logs.Methods, logLine.Method)
		}
		if !slices.Contains(Logs.Paths, logLine.Path) {
			Logs.Paths = append(Logs.Paths, logLine.Path)
		}
		if !slices.Contains(Logs.Statuses, logLine.Status) {
			Logs.Statuses = append(Logs.Statuses, logLine.Status)
		}
	}

	if bufReader.Err() != nil {
		log.Panic(bufReader.Err())
		return
	}

	var warnings int

	for _, log := range Logs.AllLogs {
		if log.Level == "WARN" {
			warnings += 1
		}
	}
	fmt.Print("There were: ", len(Logs.AllLogs), " logs. ", strconv.Itoa(warnings), " were warnings\n")
	fmt.Print(
		"There were: ", len(Logs.Methods), " Methods. ", Logs.Methods, "\n",
		"There were: ", len(Logs.Levels), " Levels. ", Logs.Levels, "\n",
		"There were: ", len(Logs.Paths), " Paths. ", Logs.Paths, "\n",
		"There were: ", len(Logs.Statuses), " Statuses. ", Logs.Statuses, "\n",
	)
}
