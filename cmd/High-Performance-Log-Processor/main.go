package main

import (
	"High-Performance-Log-Processor/internal/cli"
	"High-Performance-Log-Processor/internal/dashboard"
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ts=2026-03-14T09:01:20.006Z level=INFO req_id=0f4c9f3d method=GET path=/health status=200 duration_ms=1 bytes=2 ip=10.0.0.5 ua="kube-probe/1.31" msg="request complete"

var Logs []dashboard.LogRecord

type DateTime struct {
	Year    int
	Month   int
	Day     int
	Hour    int
	Minute  int
	Seconds int
}

func StringToInt(st string) int {
	if strings.Contains(st, ".") {
		st = st[:strings.IndexByte(st, '.')]
	}
	newInt, err := strconv.Atoi(st)
	if err != nil {
		return 0
	}
	return newInt
}

func ParseTime(inputTime string) DateTime {
	// 2026-03-14T09:01:20.006Z
	timeValue := GetValue("ts", []string{inputTime})

	splitTime := strings.Split(timeValue, "T")
	dateFields := strings.Split(splitTime[0], "-")
	timeFields := strings.Split(splitTime[1], ":")

	return DateTime{
		Year:    StringToInt(dateFields[0]),
		Month:   StringToInt(dateFields[1]),
		Day:     StringToInt(dateFields[2]),
		Hour:    StringToInt(timeFields[0]),
		Minute:  StringToInt(timeFields[1]),
		Seconds: StringToInt(timeFields[2]),
	}
}

// Take in a string
// Split on the =
func GetValue(field string, wantedKey []string) string {

	var key, value string

	for _, content := range wantedKey {
		if !strings.Contains(content, "=") {
			continue
		}
		// Split on the first = as = might be in the value but should never be in the key!
		splitLine := strings.Split(content, "=")
		// Split on a string that doesn't exist returns the string as an array with one item
		if len(splitLine) != 2 {
			continue
		}
		key = splitLine[0]
		if key == field {
			value = splitLine[1]

			if len(value) > 1 && value[0] == '"' && value[len(value)-1] == '"' {
				value = value[1 : len(value)-1]
				break
			}
		}
	}
	return value
}

// BuildLogRecord should parse a tokenized log line in a single pass and return
// the resulting LogRecord rather than repeatedly scanning for individual fields.
func BuildLogRecord(tokens []string) dashboard.LogRecord {
	// panic("not implemented")

	var key, value string
	var logLine dashboard.LogRecord

	for _, content := range tokens {
		if !strings.Contains(content, "=") {
			continue
		}
		// Split on the first = as = might be in the value but should never be in the key!
		splitLine := strings.Split(content, "=")
		// Split on a string that doesn't exist returns the string as an array with one item
		if len(splitLine) != 2 {
			continue
		}

		key = splitLine[0]
		value = splitLine[1]

		if len(value) > 1 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		switch key {
		case "ts":
			logtime, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				os.Stderr.Write([]byte("Unable to parse time! " + GetValue("ts", tokens)))
			}

			logLine.TS = logtime
		case "level":
			logLine.Level = value
		case "req_id":
			logLine.RequestID = value
		case "method":
			logLine.Method = value
		case "path":
			logLine.Path = value
		case "status":
			logLine.Status = StringToInt(value)
		case "duration_ms":
			logLine.DurationMS = StringToInt(value)
		case "bytes":
			logLine.Bytes = StringToInt(value)
		case "ip":
			logLine.IP = value
		case "ua":
			logLine.UserAgent = value
		case "msg":
			logLine.Message = value
		}
	}

	return logLine
}

func main() {

	var Flags cli.Flags

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

		Logs = append(Logs, BuildLogRecord(splitLog))
	}

	if bufReader.Err() != nil {
		log.Panic(bufReader.Err())
		return
	}

	for _, i := range dashboard.MetricsByPath(Logs) {
		err := json.NewEncoder(os.Stdout).Encode(i)
		if err != nil {
			continue
		}
	}
}
