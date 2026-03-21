package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"High-Performance-Log-Processor/internal/cli"
	"High-Performance-Log-Processor/internal/dashboard"
	"High-Performance-Log-Processor/internal/parseinput"
)

// ts=2026-03-14T09:01:20.006Z level=INFO req_id=0f4c9f3d method=GET path=/health status=200 duration_ms=1 bytes=2 ip=10.0.0.5 ua="kube-probe/1.31" msg="request complete"

func checkTimeDurationSet(inputTime time.Duration) time.Duration {

	if inputTime == 0 {
		bucket, err := time.ParseDuration("5m")
		if err != nil {
			log.Panic("[ERROR]: No time duration was found and could not be forced!")
		}
		return bucket
	}

	return 0
}

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
		log.Panic("[ERROR]: No file name found")
	} else {
		Logs = parseinput.ParseFile(Flags)
		switch {
		default:
			log.Panic("[ERROR]: No option of how to parse file")
			os.Exit(1)
		case Flags.MetricsByPath == true:
			for _, i := range dashboard.MetricsByPath(Logs) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.LatencyByPath == true:
			for _, i := range dashboard.LatencyByPath(Logs) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.SlowRequestsByPath == true:
			for _, i := range dashboard.SlowRequestsByPath(Logs) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.ErrorRateByPath == true:
			for _, i := range dashboard.ErrorRateByPath(Logs) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.RequestsByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.RequestsByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.LevelsByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.LevelsByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.WarnAndErrorCountsByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.WarnAndErrorCountsByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.StatusClassesByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.StatusClassesByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.StatusCodesByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.StatusCodesByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.MetricsByPathAndWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.MetricsByPathAndWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.SlowRequestsByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.SlowRequestsByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		case Flags.ErrorRateByWindow == true:
			Flags.Bucket = checkTimeDurationSet(Flags.Bucket)
			for _, i := range dashboard.ErrorRateByWindow(Logs, dashboard.BucketSize(Flags.Bucket)) {
				if Flags.Output == "JSON" {
					json.NewEncoder(os.Stdout).Encode(i)
				} else {
					fmt.Println(i)
				}
			}
		}

	}

}
