package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Pelegrinetti/phpublisher/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "phpublisher",
		Usage: "Composer publisher for Nexus repository.",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			cmd.PublishCmd(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
