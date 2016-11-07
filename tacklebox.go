package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {

	configData, err := RetrieveConfig()
	if err != nil {
		return
	}

	app := cli.NewApp()
	app.Name = "tacklebox"
	app.Usage = "Manage your universal assets per-project"

	app.Action = func(c *cli.Context) error {
		syncErr := configData.sync()
		if syncErr == nil {
			fmt.Println("Synced up for great good!")
		}
		return syncErr
	}

	app.Run(os.Args)
	fmt.Printf("%s", configData)
}
