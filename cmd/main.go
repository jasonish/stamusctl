package main

import (
	"os"

	cli "git.stamus-networks.com/lanath/stamus-ctl/cmd/stamusctl/commands"
	"git.stamus-networks.com/lanath/stamus-ctl/internal/app"
)

func main() {
	var err error

	switch app.Name {
	case "stamusctl":
		cli.Execute()
	default:
		cli.Execute()
	}

	if err != nil {
		os.Exit(1)
	}
}
