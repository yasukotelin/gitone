package usecase

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/micmonay/keybd_event"
)

type GitLog struct {
	No         int
	Graph      string
	CommitHash string
	Message    string
	Name       string
	Date       string
}

type GitInfo struct {
	TotalCommitCount int
	GitLogs          []GitLog
}

func GetGitInfo() (*GitInfo, error) {
	cmd := exec.Command("git", "log", "--graph", "--all", "--oneline", "--pretty=format:%d")
	graph, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	cmd = exec.Command("git", "log", "--graph", "--all", "--oneline", "--pretty=format:%h (%an) (%cd) %s", "--abbrev-commit", "--date=iso-local")
	log, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return convToGitInfo(string(graph), string(log)), nil
}

func convToGitInfo(graph, log string) *GitInfo {
	graphLines := strings.Split(graph, "\n")
	logLines := strings.Split(log, "\n")
	gitLog := make([]GitLog, 0, len(graphLines))

	// commitCount is commit number counter
	var commitCount int = 0
	for i, graphLine := range graphLines {
		if graphLine == "" {
			// Remove the last empty line.
			break
		}
		commitHash, name, date, message := parseLog(logLines[i])

		no := -1
		if commitHash != "" {
			commitCount++
			no = commitCount
		}
		gl := GitLog{
			No:         no,
			Graph:      graphLine,
			CommitHash: commitHash,
			Message:    message,
			Name:       name,
			Date:       date,
		}
		gitLog = append(gitLog, gl)
	}

	return &GitInfo{
		TotalCommitCount: commitCount,
		GitLogs:          gitLog,
	}
}

func parseLog(logLine string) (commitHash, name, date, message string) {
	// regexp is Grouping for "*    90ae84a (yasukotelin) (2020-01-21 11:30:43 +0900) Merge the message sample"
	// 1st -> commit hash
	// 2nd -> name
	// 3rd -> date
	// 4th -> message
	r := regexp.MustCompile(`([a-z0-9]*) \((.*?)\) \((.*?)\) (.*$)`)
	group := r.FindStringSubmatch(logLine)
	if len(group) != 0 {
		commitHash = group[1]
		name = group[2]
		date = group[3]
		message = group[4]
	}

	return commitHash, name, date, message
}

func RunGitShow(commitHash string) error {
	gitCmd := exec.Command("git", "show", "--color=always", commitHash)
	lessCmd := exec.Command("less", "-R")
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr
	pipe, err := gitCmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer pipe.Close()

	lessCmd.Stdin = pipe

	if err := gitCmd.Start(); err != nil {
		return err
	}
	if err := lessCmd.Start(); err != nil {
		return err
	}
	if err := sendEnterKey(); err != nil {
		return err
	}
	if err := lessCmd.Wait(); err != nil {
		return err
	}

	return nil
}

func RunGitShowStat(commitHash string) error {
	gitCmd := exec.Command("git", "show", "--color=always", "--stat", commitHash)
	lessCmd := exec.Command("less", "-R")
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr
	pipe, err := gitCmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer pipe.Close()

	lessCmd.Stdin = pipe

	if err := gitCmd.Start(); err != nil {
		return err
	}

	if err := lessCmd.Start(); err != nil {
		return err
	}
	if err := sendEnterKey(); err != nil {
		return err
	}
	if err := lessCmd.Wait(); err != nil {
		return err
	}

	return nil
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
