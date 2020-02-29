package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/yasukotelin/gitone/view"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitone"
	app.Version = "1.3.0"
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
	var theme view.Theme
	if context.String("theme") == "light" {
		theme = view.Light
	} else if context.String("theme") == "dark" {
		theme = view.Dark
	} else if context.String("theme") == "solidlight" {
		theme = view.SolidLight
	} else if context.String("theme") == "soliddark" {
		theme = view.SolidDark
	} else {
		fmt.Println("Invalid theme")
		os.Exit(1)
	}

	tui := view.NewTui(theme)
	if err := tui.Run(); err != nil {
		fmt.Println(err)
	}
}
