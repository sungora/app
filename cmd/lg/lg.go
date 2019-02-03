package main

import (
	"os"

	_ "github.com/sungora/app/lg"
	"github.com/sungora/app/startup"
)

func main() {
	os.Exit(startup.Start())
}
