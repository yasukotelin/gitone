package repo

import (
	"os"
	"os/exec"
)

func GetGitGraph() (string, error) {
	cmd := exec.Command("git", "log", "--graph", "--all", "--oneline", "--pretty=format:%d")
	graph, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(graph), err
}

func GetGitGraphWithInfo() (string, error) {
	cmd := exec.Command("git", "log", "--graph", "--all", "--oneline", "--pretty=format:%h (%an) (%cd) %s", "--abbrev-commit", "--date=iso-local")
	info, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(info), err
}

func GetGitShowWithLessCmd(commitHash string) (*exec.Cmd, error) {
	gitCmd := exec.Command("git", "show", "--color=always", commitHash)
	lessCmd := exec.Command("less", "-R")
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr
	pipe, err := gitCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	lessCmd.Stdin = pipe

	if err := gitCmd.Start(); err != nil {
		return nil, err
	}

	return lessCmd, nil
}

func GetGitShowStatWithLessCmd(commitHash string) (*exec.Cmd, error) {
	gitCmd := exec.Command("git", "show", "--color=always", "--stat", commitHash)
	lessCmd := exec.Command("less", "-R")
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr
	pipe, err := gitCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	lessCmd.Stdin = pipe

	if err := gitCmd.Start(); err != nil {
		return nil, err
	}
	return lessCmd, nil
}
