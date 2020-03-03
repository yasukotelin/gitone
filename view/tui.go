package view

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/micmonay/keybd_event"
	"github.com/rivo/tview"
	"github.com/yasukotelin/gitone/usecase"
)

type Tui struct {
	app         *tview.Application
	theme       Theme
	tuiTheme    *TuiTheme
	gitInfo     *usecase.GitInfo
	flexBox     *tview.Flex
	infoView    *tview.TextView
	treeView    *tview.List
	log1View    *tview.TextView
	log2View    *tview.TextView
	messageView *tview.TextView
}

func NewTui(theme Theme) *Tui {
	return &Tui{
		app:      tview.NewApplication(),
		theme:    theme,
		tuiTheme: NewTuiTheme(theme),
	}
}

func (t *Tui) Run() error {
	gitInfo, err := usecase.GetGitInfo()
	if err != nil {
		return err
	}
	t.gitInfo = gitInfo

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
	switch event.Key() {
	case tcell.KeyEsc:
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
	t.updateInfoView(0)
	t.updateLogView(0)
}

func (t *Tui) updateInfoView(index int) {
	t.infoView.SetText(fmt.Sprintf("(%v/%v)", t.gitInfo.GitLogs[index].No, t.gitInfo.TotalCommitCount))
}

func (t *Tui) updateLogView(index int) {
	t.log1View.SetText(getLog1Text(t.gitInfo.GitLogs[index]))
	t.log2View.SetText(getLog2Text(t.gitInfo.GitLogs[index]))
}

func getLog1Text(gitLog usecase.GitLog) string {
	if gitLog.CommitHash == "" {
		return ""
	}
	return fmt.Sprintf("[%s] %s (%s)", gitLog.CommitHash, gitLog.Name, gitLog.Date)
}

func getLog2Text(gitLog usecase.GitLog) string {
	if gitLog.CommitHash == "" {
		return ""
	}
	return gitLog.Message
}

func (t *Tui) newGitTreeView() *tview.List {
	list := tview.NewList()
	for _, log := range t.gitInfo.GitLogs {
		list.AddItem(log.Graph, "", 0, nil)
	}

	list.ShowSecondaryText(false)
	list.SetMainTextColor(t.tuiTheme.fg)
	list.SetSelectedBackgroundColor(t.tuiTheme.treeSelBgColor)
	list.SetSelectedTextColor(t.tuiTheme.treeSelFgColor)
	list.SetHighlightFullLine(true)
	list.SetWrapAround(false)
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		index := t.treeView.GetCurrentItem()
		switch event.Rune() {
		case 'j':
			t.treeView.SetCurrentItem(t.getNextCommitIdx(+1))
			return nil
		case 'k':
			t.treeView.SetCurrentItem(t.getNextCommitIdx(-1))
			return nil
		case 'g':
			return tcell.NewEventKey(tcell.KeyHome, ' ', tcell.ModNone)
		case 'G':
			return tcell.NewEventKey(tcell.KeyEnd, ' ', tcell.ModNone)
		case 's':
			{
				t.runGitShowStat(index)
				return nil
			}
		}
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			{
				t.runGitShowStat(index)
				return nil
			}
		case tcell.KeyCtrlD:
			return tcell.NewEventKey(tcell.KeyPgDn, ' ', tcell.ModNone)
		case tcell.KeyCtrlU:
			return tcell.NewEventKey(tcell.KeyPgUp, ' ', tcell.ModNone)
		case tcell.KeyDown:
			t.treeView.SetCurrentItem(t.getNextCommitIdx(+1))
			return nil
		case tcell.KeyUp:
			t.treeView.SetCurrentItem(t.getNextCommitIdx(-1))
			return nil
		}
		return event
	})
	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		t.updateInfoView(index)
		t.updateLogView(index)
	})
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortCut rune) {
		t.runGitShow(index)
	})

	list.Box.SetBackgroundColor(t.tuiTheme.bg)

	return list
}

func (t *Tui) getNextCommitIdx(direction int) int {
	orgIdx := t.treeView.GetCurrentItem()
	idx := orgIdx
	for {
		idx = idx + direction
		if idx < 0 || idx == t.treeView.GetItemCount() {
			return orgIdx
		}
		if t.gitInfo.GitLogs[idx].CommitHash == "" {
			continue
		} else {
			return idx
		}
	}
}

func (t *Tui) runGitShow(index int) {
	commitHash := t.gitInfo.GitLogs[index].CommitHash
	if commitHash == "" {
		return
	}
	t.app.Suspend(func() {
		cmd, err := usecase.GetGitShowWithLessCmd(commitHash)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := sendEnterKey(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})
}

func (t *Tui) runGitShowStat(index int) {
	commitHash := t.gitInfo.GitLogs[index].CommitHash
	if commitHash == "" {
		return
	}
	t.app.Suspend(func() {
		cmd, err := usecase.GetGitShowStatWithLessCmd(commitHash)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := sendEnterKey(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})
}

func (t *Tui) newInfoView() *tview.TextView {
	return newTextView(t.tuiTheme.fg, t.tuiTheme.bg)
}

func (t *Tui) newLog1View() *tview.TextView {
	return newTextView(t.tuiTheme.log1Fg, t.tuiTheme.log1Bg)
}

func (t *Tui) newLog2View() *tview.TextView {
	return newTextView(t.tuiTheme.fg, t.tuiTheme.bg)
}

func (t *Tui) newMessageView() *tview.TextView {
	return newTextView(tcell.ColorRed, t.tuiTheme.bg)
}

func newTextView(textColor, bgColor tcell.Color) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextColor(textColor)
	tv.SetBackgroundColor(bgColor)
	return tv
}

// SendEnterKey have to be used by on the suspend function.
// Because tcell has bug https://github.com/gdamore/tcell/issues/194,
// suspend function will lost first key.
func sendEnterKey() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	kb.SetKeys(keybd_event.VK_ENTER)
	return kb.Launching()
}
