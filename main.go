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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainAction(context *cli.Context) {
	gui := NewGui(Light)
	if err := gui.Run(); err != nil {
		fmt.Println(err)
	}
}
