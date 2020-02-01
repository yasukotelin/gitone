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
	app.Version = "1.1.0"
	app.Usage = "gitone is simple git tree viewer"
	app.Action = mainAction
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "theme",
			Value: "dark",
			Usage: "dark/light/soliddark/solidlight",
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
	} else if context.String("theme") == "solidlight" {
		theme = SolidLight
	} else if context.String("theme") == "soliddark" {
		theme = SolidDark
	} else {
		fmt.Println("Invalid theme")
		os.Exit(1)
	}

	tui := NewTui(theme)
	if err := tui.Run(); err != nil {
		fmt.Println(err)
	}
}
