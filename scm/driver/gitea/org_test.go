// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestOrgFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/orgs/gogits").
		Reply(200).
		Type("application/json").
		File("testdata/organization.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Organizations.Find(context.Background(), "gogits")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Organization)
	raw, _ := os.ReadFile("testdata/organization.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrgList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/user/orgs").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/organizations.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Organizations.List(context.Background(), &scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Organization{}
	raw, _ := os.ReadFile("testdata/organizations.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}
