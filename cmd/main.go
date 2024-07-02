package main

import (

	// cli "stamus-ctl/cmd/stamusctl/commands"
	"os"
	ctl "stamus-ctl/cmd/ctl"
	daemon "stamus-ctl/cmd/daemon"
	"stamus-ctl/internal/app"
)

func main() {
	// var err error

	switch app.Name {
	case "stamusctl":
		ctl.Execute()
	case "stamusd":
		daemon.Execute()
	default:
		exec := os.Getenv("STAMUS_APP_NAME")
		switch exec {
		case "stamusctl":
			ctl.Execute()
		case "stamusd":
			daemon.Execute()
		default:
			daemon.Execute()
		}
	}

	// if err != nil {
	// 	os.Exit(1)
	// }
}
