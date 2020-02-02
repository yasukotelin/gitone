package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConvToGitLogs(t *testing.T) {
	gf, err := os.Open("./_testdata/graph.txt")
	if err != nil {
		t.Error(err)
	}
	defer gf.Close()
	lf, err := os.Open("./_testdata/log.txt")
	if err != nil {
		t.Error(err)
	}

	gb, err := ioutil.ReadAll(gf)
	lb, err := ioutil.ReadAll(lf)
	gitInfo := convToGitInfo(string(gb), string(lb))

	// test the length
	if len(gitInfo.GitLogs) != 14 {
		t.Errorf("Git log length should be 14 but actual is %v", len(gitInfo.GitLogs))
	}
	// test the total commit number
	if gitInfo.TotalCommitCount != 11 {
		t.Errorf("Total commit number should be 11 but actual is %v", gitInfo.TotalCommitCount)
	}
	// test to extract the graph
	if gitInfo.GitLogs[0].Graph != "*  (HEAD -> master, tag: v2.19.1, origin/master, origin/HEAD)" {
		t.Error("Failed the extract the graph")
	}
	// test to extract the No
	if gitInfo.GitLogs[5].No != 5 {
		t.Errorf("Failed the extract the No. extract is 4 but actual is %v", gitInfo.GitLogs[5].No)
	}
	// test to extract the commit hash
	if gitInfo.GitLogs[2].CommitHash != "90ae84a" {
		t.Error("Failed the extract the commit hash")
	}
	if gitInfo.GitLogs[12].CommitHash != "2b30b7412413" {
		t.Error("Failed the extract the commit hash")
	}
	// test to extract the message
	if gitInfo.GitLogs[5].Message != "Add message sample. and function() sample." {
		t.Error("Failed the extract the message")
	}
	// test to extract the name
	if gitInfo.GitLogs[7].Name != "yasukotelin" {
		t.Error("Failed the extract the name")
	}
	// test to extract the date.
	if gitInfo.GitLogs[13].Date != "2019-06-26 03:17:17 +0900" {
		t.Error("Failed the extract the date")
	}
}
