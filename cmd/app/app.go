package main

import (
	"os"

	_ "github.com/sungora/app/cmd/workflow/worker"
	_ "github.com/sungora/app/core"
	_ "github.com/sungora/app/server"
	"github.com/sungora/app/startup"
)

func main() {
	os.Exit(startup.Start())
}
