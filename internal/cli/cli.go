package cli

import (
	"fmt"
	"log"
	"os"
	"time"
)

const version string = "0.2.0"

type Flags struct {
	FileName                   string
	Ping                       bool
	Bucket                     time.Duration
	Output                     string
	MetricsByPath              bool
	LatencyByPath              bool
	SlowRequestsByPath         bool
	ErrorRateByPath            bool
	RequestsByWindow           bool
	LevelsByWindow             bool
	WarnAndErrorCountsByWindow bool
	StatusClassesByWindow      bool
	StatusCodesByWindow        bool
	MetricsByPathAndWindow     bool
	SlowRequestsByWindow       bool
	ErrorRateByWindow          bool
}

func CLI(args []string) Flags {
	var returnFlags Flags

	for index := 0; index < len(args); index++ {
		arg := args[index]

		switch arg {
		default:
			if _, err := os.Lstat(arg); err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR]: did not recognise the command: %s\n", arg)
			}
			returnFlags.FileName = arg
		case "help":
			fmt.Fprint(os.Stdout, usageText)
			return Flags{}
		case "version":
			fmt.Fprintf(os.Stdout, "High-Performance-Log-Processor %s\n", version)
			return Flags{}
		case "output", "--output", "-o", "-output":
			if index+1 < len(args) {
				if args[index+1] == "JSON" || args[index+1] == "json" || args[index+1] == "Json" || args[index+1] == "JSon" || args[index+1] == "JSOn" {
					returnFlags.Output = "JSON"
				}
				index += 1
			} else {
				log.Print("[WARNING]: no output format detected")
			}
		case "ping":
			return Flags{
				Ping: true,
			}
		case "--time":
			if index+1 < len(args) {
				bucket, err := time.ParseDuration(args[index+1])
				if err != nil {
					log.Fatal("Unable to parse time: ", err)
				}
				returnFlags.Bucket = bucket
				index++
			} else {
				log.Print("[WARNING]: no time period entered, defaulting to 5 minutes")
				bucket, err := time.ParseDuration("5m")
				if err != nil {
					log.Fatal("Unable to parse time: ", err)
				}
				returnFlags.Bucket = bucket
			}
		case "MetricsByPath":
			returnFlags.MetricsByPath = true
		case "LatencyByPath":
			returnFlags.LatencyByPath = true
		case "SlowRequestsByPath":
			returnFlags.SlowRequestsByPath = true
		case "ErrorRateByPath":
			returnFlags.ErrorRateByPath = true
		case "RequestsByWindow":
			returnFlags.RequestsByWindow = true
		case "LevelsByWindow":
			returnFlags.LevelsByWindow = true
		case "WarnAndErrorCountsByWindow":
			returnFlags.WarnAndErrorCountsByWindow = true
		case "StatusClassesByWindow":
			returnFlags.StatusClassesByWindow = true
		case "StatusCodesByWindow":
			returnFlags.StatusCodesByWindow = true
		case "MetricsByPathAndWindow":
			returnFlags.MetricsByPathAndWindow = true
		case "SlowRequestsByWindow":
			returnFlags.SlowRequestsByWindow = true
		case "ErrorRateByWindow":
			returnFlags.ErrorRateByWindow = true
		}
	}

	return returnFlags
}
