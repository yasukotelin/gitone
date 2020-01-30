package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitone"
	app.Version = "1.0.0"
	app.Usage = "gitone is simple git tree viewer"
	app.Action = mainAction
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "theme",
			Value: "dark",
			Usage: "theme is dark or light",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainAction(context *cli.Context) {

	var theme Theme
	if context.String("theme") == "light" {
		theme = Light
	} else if context.String("theme") == "dark" {
		theme = Dark
	} else {
		fmt.Println("--theme should be dark or light")
		os.Exit(1)
	}

	tui := NewTui(theme)
	if err := tui.Run(); err != nil {
		fmt.Println(err)
	}
}
