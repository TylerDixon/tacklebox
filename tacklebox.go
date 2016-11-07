package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"path"
)

func main() {

	configData, err := RetrieveConfig()
	if err != nil {
		return
	}

	app := cli.NewApp()
	app.Name = "tacklebox"
	app.Usage = "Manage your universal assets per-project"
	app.Commands = []cli.Command{
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Sync files declared in configurations",
			Action: func(c *cli.Context) error {
				syncErr := configData.Sync()
				if syncErr == nil {
					fmt.Println("Synced up for great good!")
				}
				return syncErr
			},
		},
		{
			Name:    "readdir",
			Aliases: []string{"r"},
			Usage:   "Add all directories in given directory to config Projects for configuration",
			Action: func(c *cli.Context) error {
				dir := "."
				if c.NArg() > 0 {
					dir = c.Args().Get(0)
				}
				cwd, getWdErr := os.Getwd()
				if getWdErr != nil {
					fmt.Printf("Failed to get cwd due to error %s", getWdErr)
					return getWdErr
				}
				configDirsErr := configData.ConfigDirs(path.Join(cwd, dir))
				if configDirsErr != nil {
					fmt.Printf("Failed to get nested directories due to error %s", configDirsErr)
					return configDirsErr
				}
				saveErr := configData.Save()
				if saveErr != nil {
					fmt.Printf("Failed to get nested directories due to error %s", saveErr)
					return saveErr
				}
				return nil

			},
		},
	}
	app.Run(os.Args)
	fmt.Printf("%s", configData)
}
