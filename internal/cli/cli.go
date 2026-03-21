package cli

import (
	"log"
	"os"
	"time"
)

const version string = "0.0.2"

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
				os.Stderr.Write([]byte("[ERROR]: Did not recognise the command; " + arg))
			}
			returnFlags.FileName = arg
		case "help":
			os.Stderr.Write([]byte("USEAGE: pass in a file name"))
		case "version":
			os.Stderr.Write([]byte("Version: " + version))
		case "output", "--output", "-o", "-output":
			if len(args) >= index+1 {

			} else {
				log.Panic("[ERROR]: No output detected")
			}
		case "ping":
			return Flags{
				Ping: true,
			}
		case "--time":
			if len(args) >= index+1 {
				bucket, err := time.ParseDuration(args[index+1])
				if err != nil {
					log.Fatal("Unable to parse time: ", err)
				}
				returnFlags.Bucket = bucket
				index++
			} else {
				log.Print("[waRNING]: No time period entered, defaulting to 5 minutes")
				bucket, err := time.ParseDuration("5m")
				if err != nil {
					log.Fatal("Unable to parse time: ", err)
				}
				returnFlags.Bucket = bucket
			}
			index += 1
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
