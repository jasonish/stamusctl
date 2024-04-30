package main

import (

	// cli "stamus-ctl/cmd/stamusctl/commands"
	ctl "stamus-ctl/cmd/ctl"
	"stamus-ctl/internal/app"
)

func main() {
	// var err error

	switch app.Name {
	case "stamusctl":
		ctl.Execute()
	default:
		ctl.Execute()
	}

	// if err != nil {
	// 	os.Exit(1)
	// }
}
