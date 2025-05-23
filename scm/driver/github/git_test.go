// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit.json")

	client := NewDefault()
	got, res, err := client.Git.FindCommit(context.Background(), "octocat/hello-world", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := os.ReadFile("testdata/commit.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/branches/master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch.json")

	client := NewDefault()
	got, res, err := client.Git.FindBranch(context.Background(), "octocat/hello-world", "master")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/branch.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindTag(t *testing.T) {
	git := new(gitService)
	_, _, err := git.FindTag(context.Background(), "octocat/hello-world", "v1.0")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/commits").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("ref", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/commits.json")

	client := NewDefault()
	got, res, err := client.Git.ListCommits(context.Background(), "octocat/hello-world", scm.CommitListOptions{Ref: "master", Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := os.ReadFile("testdata/commits.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/branches").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/branches.json")

	client := NewDefault()
	got, res, err := client.Git.ListBranches(context.Background(), "octocat/hello-world", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/branches.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/tags").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/tags.json")

	client := NewDefault()
	got, res, err := client.Git.ListTags(context.Background(), "octocat/hello-world", &scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := os.ReadFile("testdata/tags.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/changes.json")

	client := NewDefault()
	got, res, err := client.Git.ListChanges(context.Background(), "octocat/hello-world", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/changes.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCompareCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/compare/762941318ee16e59dabbacb1b4049eec22f0d303...7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/changes.json")

	client := NewDefault()
	got, res, err := client.Git.CompareCommits(context.Background(), "octocat/hello-world", "762941318ee16e59dabbacb1b4049eec22f0d303", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", &scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := os.ReadFile("testdata/changes.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCreateRef(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/octocat/hello-world/git/refs").
		Reply(http.StatusCreated).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/ref.json")

	client := NewDefault()
	got, res, err := client.Git.CreateRef(context.Background(), "octocat/hello-world", "refs/heads/testing", "aa218f56b14c9653891f9e74264a383fa43fefbd")
	if err != nil {
		t.Error(err)
		return
	}

	want := &scm.Reference{}
	raw, _ := os.ReadFile("testdata/ref.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitGetDefaultBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/repo.json")

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/branches/master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch.json")

	client := NewDefault()
	got, res, err := client.Git.GetDefaultBranch(context.Background(), "octocat/hello-world")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := os.ReadFile("testdata/branch.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
