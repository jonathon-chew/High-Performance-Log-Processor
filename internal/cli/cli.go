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
				// fmt.Fprint(os.Stderr, usageText)
			}
			returnFlags.FileName = arg
		case "help", "--help", "-h", "--useage", "-u":
			fmt.Fprint(os.Stdout, usageText)
			os.Exit(0)
		case "version":
			fmt.Fprintf(os.Stdout, "High-Performance-Log-Processor %s\n", version)
			os.Exit(0)
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
		case "MetricsByPath", "--MetricsByPath", "-MetricsByPath", "MBP", "-MBP", "--MBP":
			returnFlags.MetricsByPath = true
		case "LatencyByPath", "--LatencyByPath", "-LatencyByPath", "LBP", "-LBP", "--LBP":
			returnFlags.LatencyByPath = true
		case "SlowRequestsByPath", "--SlowRequestsByPath", "-SlowRequestsByPath", "SBP", "-SBP", "--SBP":
			returnFlags.SlowRequestsByPath = true
		case "ErrorRateByPath", "--ErrorRateByPath", "-ErrorRateByPath", "EBP", "-EBP", "--EBP":
			returnFlags.ErrorRateByPath = true
		case "RequestsByWindow", "--RequestsByWindow", "-RequestsByWindow", "RBW", "-RBW", "--RBW":
			returnFlags.RequestsByWindow = true
		case "LevelsByWindow", "--LevelsByWindow", "-LevelsByWindow", "LBW", "-LBW", "--LBW":
			returnFlags.LevelsByWindow = true
		case "WarnAndErrorCountsByWindow", "--WarnAndErrorCountsByWindow", "-WarnAndErrorCountsByWindow", "WBW", "-WBW", "--WBW":
			returnFlags.WarnAndErrorCountsByWindow = true
		case "StatusClassesByWindow", "--StatusClassesByWindow", "-StatusClassesByWindow", "SClBW", "--SClBW", "-SClBW":
			returnFlags.StatusClassesByWindow = true
		case "StatusCodesByWindow", "--StatusCodesByWindow", "-StatusCodesByWindow", "SCoBW", "-SCoBW", "--SCoBW":
			returnFlags.StatusCodesByWindow = true
		case "MetricsByPathAndWindow", "--MetricsByPathAndWindow", "-MetricsByPathAndWindow", "MBPAW", "--MBPAW", "-MBPAW":
			returnFlags.MetricsByPathAndWindow = true
		case "SlowRequestsByWindow", "--SlowRequestsByWindow", "-SlowRequestsByWindow", "SBW", "-SBW", "--SBW":
			returnFlags.SlowRequestsByWindow = true
		case "ErrorRateByWindow", "--ErrorRateByWindow", "-ErrorRateByWindow", "EBW", "-EBW", "--EBW":
			returnFlags.ErrorRateByWindow = true
		}
	}

	return returnFlags
}
