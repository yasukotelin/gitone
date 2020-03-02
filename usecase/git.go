package usecase

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/yasukotelin/gitone/repo"
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
	graph, err := repo.GetGitGraph()
	if err != nil {
		return nil, err
	}
	info, err := repo.GetGitGraphWithInfo()
	if err != nil {
		return nil, err
	}

	return convToGitInfo(graph, info), nil
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

func GetGitShowWithLessCmd(commitHash string) (*exec.Cmd, error) {
	return repo.GetGitShowWithLessCmd(commitHash)
}

func GetGitShowStatWithLessCmd(commitHash string) (*exec.Cmd, error) {
	return repo.GetGitShowStatWithLessCmd(commitHash)
}
