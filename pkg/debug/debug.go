package debug

import "log"

// Enabling this will enable debug logging
var Enabled = false

func Log(msg string) {
	if Enabled {
		log.Printf(msg)
	}
}

func Logf(msg string, args ...interface{}) {
	if Enabled {
		log.Printf(msg, args...)
	}
}

func LogDump(msg []byte, err error) {
	if Enabled {
		if err != nil {
			panic(err)
		}

		log.Printf(string(msg))
	}
}
