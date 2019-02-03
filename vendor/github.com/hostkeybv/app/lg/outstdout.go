package lg

import (
	"fmt"
	"os"
)

func saveStdout(m msg) {
	switch m.Level {
	case "fatal":
		fmt.Fprintf(os.Stderr, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	case "critical":
		fmt.Fprintf(os.Stderr, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	case "error":
		fmt.Fprintf(os.Stderr, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	case "warning":
		fmt.Fprintf(os.Stdout, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	case "notice":
		fmt.Fprintf(os.Stdout, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	case "info":
		fmt.Fprintf(os.Stdout, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	case "debug":
		fmt.Fprintf(os.Stdout, "%s [%s] %s %d [%s]\n", m.Datetime, m.Level, m.Message, m.LineNumber, m.Action)
	}
}
