package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Theme int

const (
	Dark Theme = iota
	Light
)

type Gui struct {
	app   *tview.Application
	theme Theme
	bg    tcell.Color
	fg    tcell.Color
}

func NewGui(theme Theme) *Gui {
	return &Gui{
		app:   tview.NewApplication(),
		theme: theme,
		bg:    getBg(theme),
		fg:    getFg(theme),
	}
}

func getBg(theme Theme) tcell.Color {
	if theme == Dark {
		return tcell.ColorBlack
	} else {
		return tcell.ColorWhite
	}
}

func getFg(theme Theme) tcell.Color {
	if theme == Dark {
		return tcell.ColorWhite
	} else {
		return tcell.ColorBlack
	}
}

func (g *Gui) Run() error {
	logs, err := getGitLog()
	if err != nil {
		return err
	}

	list := tview.NewList()
	for _, log := range logs {
		list.AddItem(log.Graph, "", 0, nil)
	}

	list.ShowSecondaryText(false)
	list.SetMainTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorBlue)
	list.SetSelectedTextColor(tcell.ColorWhite)
	list.SetHighlightFullLine(true)
	list.SetWrapAround(false)

	list.Box.SetBackgroundColor(tcell.ColorWhite)

	if err := g.app.SetRoot(list, true).Run(); err != nil {
		return err
	}

	return nil
}
