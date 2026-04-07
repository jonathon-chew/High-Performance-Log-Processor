package parseinput

import (
	"bufio"
	"log"
	"os"

	"github.com/jonathon-chew/High-Performance-Log-Processor/internal/cli"
	"github.com/jonathon-chew/High-Performance-Log-Processor/internal/dashboard"
)

func ParseFile(flags cli.Flags) []dashboard.LogRecord {
	filePointer, err := os.Open(flags.FileName)
	if err != nil {
		log.Print("[ERROR]: Unable to open file ", err)
		return []dashboard.LogRecord{}
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
		return []dashboard.LogRecord{}
	}

	return Logs
}
