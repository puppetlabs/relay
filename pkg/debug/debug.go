package debug

import "log"
import "io"
import "os"

var logger = log.New(os.Stdout, "[debug] ", log.LstdFlags)

// Enabling this will enable debug logging
var Enabled = false

func Log(msg string) {
	if Enabled {
		logger.Printf(msg)
	}
}

func Logf(msg string, args ...interface{}) {
	if Enabled {
		logger.Printf(msg, args...)
	}
}

func LogDump(msg []byte, err error) {
	if Enabled {
		if err != nil {
			panic(err)
		}

		logger.Printf(string(msg))
	}
}

type noopWriter struct{}

func (noopWriter) Write(buf []byte) (int, error) { return len(buf), nil }

func Writer() io.Writer {
	if Enabled {
		return logger.Writer()
	} else {
		return noopWriter{}
	}
}
