package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Theme int

const (
	Dark Theme = iota
	Light
	SolidDark
	SolidLight
)

type Tui struct {
	app         *tview.Application
	theme       Theme
	color       *TuiColor
	gitLogs     []GitLog
	flexBox     *tview.Flex
	infoView    *tview.TextView
	treeView    *tview.List
	log1View    *tview.TextView
	log2View    *tview.TextView
	messageView *tview.TextView
}

type TuiColor struct {
	bg             tcell.Color
	fg             tcell.Color
	treeSelFgColor tcell.Color
	treeSelBgColor tcell.Color
	log1Bg         tcell.Color
	log1Fg         tcell.Color
}

func NewTui(theme Theme) *Tui {
	if theme == Dark {
		return NewDarkTui()
	} else if theme == Light {
		return NewLightTui()
	} else if theme == SolidDark {
		return NewSolidDarkTui()
	} else {
		return NewSolidLightTui()
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
		},
	}
}

func NewSolidDarkTui() *Tui {
	return &Tui{
		app:   tview.NewApplication(),
		theme: Dark,
		color: &TuiColor{
			bg:             tcell.ColorBlack,
			fg:             tcell.ColorWhite,
			treeSelBgColor: tcell.ColorWhite,
			treeSelFgColor: tcell.ColorBlack,
			log1Bg:         tcell.ColorWhite,
			log1Fg:         tcell.ColorBlack,
		},
	}
}
func NewSolidLightTui() *Tui {
	return &Tui{
		app:   tview.NewApplication(),
		theme: Dark,
		color: &TuiColor{
			bg:             tcell.ColorWhite,
			fg:             tcell.ColorBlack,
			treeSelBgColor: tcell.ColorBlack,
			treeSelFgColor: tcell.ColorWhite,
			log1Bg:         tcell.ColorBlack,
			log1Fg:         tcell.ColorWhite,
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
	t.infoView = t.newInfoView()
	t.treeView = t.newGitTreeView()
	t.log1View = t.newLog1View()
	t.log2View = t.newLog2View()
	t.messageView = t.newMessageView()

	t.flexBox = tview.NewFlex().SetDirection(tview.FlexRow)
	t.flexBox.AddItem(t.infoView, 1, 1, false)
	t.flexBox.AddItem(t.treeView, 0, 1, true)
	t.flexBox.AddItem(t.log1View, 1, 1, false)
	t.flexBox.AddItem(t.log2View, 1, 1, false)
	t.flexBox.AddItem(t.messageView, 1, 1, false)

	// Init show
	t.updateInfoView()
	t.updateLogView(t.gitLogs[0])

	// Init event
	t.treeView.SetChangedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		t.updateInfoView()
		t.updateLogView(t.gitLogs[index])
	})
	t.treeView.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		commitHash := t.gitLogs[index].CommitHash
		if commitHash == "" {
			return
		}
		t.app.Suspend(func() {
			if err := RunGitShow(commitHash); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		})
	})
}

func (t *Tui) updateInfoView() {
	total := len(t.gitLogs)
	current := t.treeView.GetCurrentItem() + 1
	t.infoView.SetText(fmt.Sprintf("(%v/%v)", current, total))
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
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
		case 'g':
			return tcell.NewEventKey(tcell.KeyHome, ' ', tcell.ModNone)
		case 'G':
			return tcell.NewEventKey(tcell.KeyEnd, ' ', tcell.ModNone)
		}
		switch event.Key() {
		case tcell.KeyCtrlD:
			return tcell.NewEventKey(tcell.KeyPgDn, ' ', tcell.ModNone)
		case tcell.KeyCtrlU:
			return tcell.NewEventKey(tcell.KeyPgUp, ' ', tcell.ModNone)
		}
		return event
	})

	list.Box.SetBackgroundColor(t.color.bg)

	return list
}

func (t *Tui) newInfoView() *tview.TextView {
	return newView(t.color.fg, t.color.bg)
}

func (t *Tui) newLog1View() *tview.TextView {
	return newView(t.color.log1Fg, t.color.log1Bg)
}

func (t *Tui) newLog2View() *tview.TextView {
	return newView(t.color.fg, t.color.bg)
}

func (t *Tui) newMessageView() *tview.TextView {
	return newView(tcell.ColorRed, t.color.bg)
}

func newView(textColor, bgColor tcell.Color) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextColor(textColor)
	tv.SetBackgroundColor(bgColor)
	return tv
}
