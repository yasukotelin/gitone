package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Theme int

const (
	Dark Theme = iota
	Light
)

type Tui struct {
	app      *tview.Application
	theme    Theme
	color    *TuiColor
	gitLogs  []GitLog
	flexBox  *tview.Flex
	treeView *tview.List
	log1View *tview.TextView
	log2View *tview.TextView
}

type TuiColor struct {
	bg             tcell.Color
	fg             tcell.Color
	treeSelFgColor tcell.Color
	treeSelBgColor tcell.Color
	log1Bg         tcell.Color
	log1Fg         tcell.Color
	log2Bg         tcell.Color
	log2Fg         tcell.Color
}

func NewTui(theme Theme) *Tui {
	if theme == Dark {
		return NewDarkTui()
	} else {
		return NewLightTui()
	}
}

func NewLightTui() *Tui {
	return &Tui{
		app:   tview.NewApplication(),
		theme: Light,
		color: &TuiColor{
			bg:             tcell.ColorWhite,
			fg:             tcell.ColorBlack,
			treeSelBgColor: tcell.ColorGreen,
			treeSelFgColor: tcell.ColorWhite,
			log1Bg:         tcell.ColorBlue,
			log1Fg:         tcell.ColorWhite,
			log2Bg:         tcell.ColorWhite,
			log2Fg:         tcell.ColorBlack,
		},
	}
}

func NewDarkTui() *Tui {
	return &Tui{
		app:   tview.NewApplication(),
		theme: Dark,
		color: &TuiColor{
			bg:             tcell.ColorBlack,
			fg:             tcell.ColorWhite,
			treeSelBgColor: tcell.ColorGreen,
			treeSelFgColor: tcell.ColorWhite,
			log1Bg:         tcell.ColorBlue,
			log1Fg:         tcell.ColorWhite,
			log2Bg:         tcell.ColorBlack,
			log2Fg:         tcell.ColorWhite,
		},
	}
}

func (t *Tui) Run() error {
	gitLogs, err := GetGitLog()
	if err != nil {
		return err
	}
	t.gitLogs = gitLogs

	t.initView()

	t.app.SetInputCapture(t.inputCapture)
	if err := t.app.SetRoot(t.flexBox, true).Run(); err != nil {
		return err
	}

	return nil
}

func (t *Tui) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		t.app.Stop()
	}
	return event
}

func (t *Tui) initView() {
	t.treeView = t.newGitTreeView()
	t.log1View = t.newLog1View()
	t.log2View = t.newLog2View()

	t.flexBox = tview.NewFlex().SetDirection(tview.FlexRow)
	t.flexBox.AddItem(t.treeView, 0, 1, true)
	t.flexBox.AddItem(t.log1View, 1, 1, false)
	t.flexBox.AddItem(t.log2View, 1, 1, false)

	// Init show
	t.updateLogView(t.gitLogs[0])

	// Init event
	t.treeView.SetChangedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		t.updateLogView(t.gitLogs[index])
	})
	t.treeView.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		t.app.Suspend(func() {
			RunGitShow(t.gitLogs[index].CommitHash)
		})
	})
}

func (t *Tui) updateLogView(gitLog GitLog) {
	t.log1View.SetText(getLog1Text(gitLog))
	t.log2View.SetText(getLog2Text(gitLog))
}

func getLog1Text(gitLog GitLog) string {
	if gitLog.CommitHash == "" {
		return ""
	}
	return fmt.Sprintf("[%s] %s (%s)", gitLog.CommitHash, gitLog.Name, gitLog.Date)
}

func getLog2Text(gitLog GitLog) string {
	if gitLog.CommitHash == "" {
		return ""
	}
	return gitLog.Message
}

func (t *Tui) newGitTreeView() *tview.List {
	list := tview.NewList()
	for _, log := range t.gitLogs {
		list.AddItem(log.Graph, "", 0, nil)
	}

	list.ShowSecondaryText(false)
	list.SetMainTextColor(t.color.fg)
	list.SetSelectedBackgroundColor(t.color.treeSelBgColor)
	list.SetSelectedTextColor(t.color.treeSelFgColor)
	list.SetHighlightFullLine(true)
	list.SetWrapAround(false)

	list.Box.SetBackgroundColor(t.color.bg)

	return list
}

func (t *Tui) newLog1View() *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextColor(t.color.log1Fg)
	tv.SetBackgroundColor(t.color.log1Bg)
	return tv
}

func (t *Tui) newLog2View() *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextColor(t.color.log2Fg)
	tv.SetBackgroundColor(t.color.log2Bg)
	return tv
}
