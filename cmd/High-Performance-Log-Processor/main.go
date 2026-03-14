package main

import (
	"High-Performance-Log-Processor/internal/cli"
	"High-Performance-Log-Processor/internal/dashboard"
	"bufio"
	"fmt"
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
	timeValue := GetValue(inputTime)

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

		logtime, err := time.Parse(time.RFC3339Nano, GetValue(splitLog[0]))
		if err != nil {
			log.Println("Unable to parse time!")
		}

		var logLine = dashboard.LogRecord{
			TS:         logtime,
			Level:      GetValue(splitLog[1]),
			RequestID:  GetValue(splitLog[2]),
			Method:     GetValue(splitLog[3]),
			Path:       GetValue(splitLog[4]),
			Status:     StringToInt(GetValue(splitLog[5])),
			DurationMS: StringToInt(GetValue(splitLog[6])),
			Bytes:      StringToInt(GetValue(splitLog[7])),
			IP:         GetValue(splitLog[8]),
			// Hint: quoted fields like ua and msg still include their surrounding
			// quotes at this stage. Trimming those is a separate parsing step.
			UserAgent: GetValue(splitLog[9]),
			Message:   GetValue(splitLog[10]),
		}

		Logs = append(Logs, logLine)
	}

	if bufReader.Err() != nil {
		log.Panic(bufReader.Err())
		return
	}

	for _, i := range dashboard.MetricsByPath(Logs) {
		fmt.Println(i)
	}

	/* fmt.Print("There were: ", len(Logs.AllLogs), " logs. ", strconv.Itoa(warnings), " were warnings\n")
	fmt.Print(
		"There were: ", len(Logs.Methods), " Methods. ", Logs.Methods, "\n",
		"There were: ", len(Logs.Levels), " Levels. ", Logs.Levels, "\n",
		"There were: ", len(Logs.Paths), " Paths. ", Logs.Paths, "\n",
		"There were: ", len(Logs.Statuses), " Statuses. ", Logs.Statuses, "\n",
	) */
}
