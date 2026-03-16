package ping

import (
	"High-Performance-Log-Processor/internal/dashboard"
	"High-Performance-Log-Processor/internal/parseinput"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func ParsePing() {
	scanner := bufio.NewScanner(os.Stdin)
	lastRowCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		var ip string

		if strings.Index(line, ":") != -1 {
			ip = line[:strings.Index(line, ":")]
		} else {
			ip = "Failed ping"
		}

		parts := strings.Split(line, " ")
		var path string

		if path = parseinput.GetValue("ttl", parts); path == "" {
			path = "Error Path"
		}

		// 64 bytes from 8.8.8.8: icmp_seq=1968 ttl=117 time=64.275 ms
		var templog = dashboard.LogRecord{
			TS:         time.Now(),
			IP:         ip,
			DurationMS: parseinput.StringToInt(parseinput.GetValue("time", parts)),
			Path:       path,
		}

		parseinput.Logs = append(parseinput.Logs, templog)
		/* if err := json.NewEncoder(os.Stdout).Encode(dashboard.MetricsByPath(Logs)); err != nil {
			continue
		} */

		message, err := json.Marshal(dashboard.MetricsByPath(parseinput.Logs))
		if err != nil {
			continue
		}
		fmt.Println(string(message))
	}

	if lastRowCount > 0 {
		fmt.Print("\n")
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
}
