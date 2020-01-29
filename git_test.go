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
	gitLogs := convToGitLogs(string(gb), string(lb))

	// test the length
	if len(gitLogs) != 14 {
		t.Errorf("Git log length should be 14 but actual is %v", len(gitLogs))
	}
	// test to extract the graph
	if gitLogs[0].Graph != "*  (HEAD -> master, tag: v2.19.1, origin/master, origin/HEAD)" {
		t.Error("Failed the extract the graph")
	}
	// test to extract the commit hash
	if gitLogs[2].CommitHash != "90ae84a" {
		t.Error("Failed the extract the commit hash")
	}
	// test to extract the message
	if gitLogs[1].Message != "Update a message sample" {
		t.Error("Failed the extract the message")
	}
	// test to extract the name
	if gitLogs[7].Name != "yasukotelin" {
		t.Error("Failed the extract the name")
	}
	// test to extract the date.
	if gitLogs[13].Date != "2019-06-26 03:17:17 +0900" {
		t.Error("Failed the extract the date")
	}
}
