package main

import (
	"os"

	_ "github.com/sungora/app/cmd/app/config"
	_ "github.com/sungora/app/cmd/app/controller"
	_ "github.com/sungora/app/cmd/app/model"
	_ "github.com/sungora/app/cmd/app/worker"
	_ "github.com/sungora/app/server"

	"github.com/sungora/app/core"
)

func main() {
	os.Exit(core.Start())
}
