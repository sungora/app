package lg

import (
	"fmt"
	"os"
)

func saveStdout(m msg) {
	switch m.Level {
	case "fatal":
		fmt.Fprintf(os.Stderr, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	case "critical":
		fmt.Fprintf(os.Stderr, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	case "error":
		fmt.Fprintf(os.Stderr, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	case "warning":
		fmt.Fprintf(os.Stdout, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	case "notice":
		fmt.Fprintf(os.Stdout, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	case "info":
		fmt.Fprintf(os.Stdout, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	case "debug":
		fmt.Fprintf(os.Stdout, "%s [%s] %s [%s:%d]\n", m.Datetime, m.Level, m.Message, m.Action, m.LineNumber)
	}
}
