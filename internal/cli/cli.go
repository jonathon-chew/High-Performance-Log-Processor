package cli

import (
	"os"
)

const version string = "0.0.1"

type Flags struct {
	FileName string
	Ping     bool
}

func CLI(args []string) Flags {

	var returnFlags Flags

	for _, arg := range args {
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
		case "ping":
			return Flags{
				Ping: true,
			}
		}
	}

	return returnFlags
}
