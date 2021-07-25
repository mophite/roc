// Copyright (c) 2021 roc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

// Package rocx enable project default generator
package rocx

import (
    "sort"
    "strings"

    "github.com/go-roc/roc/rlog"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/config"
    "gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitInfo struct {
    latestTag     string
    enableDynamic bool
    remote        *git.Remote
}

var defaultLatestTag = "v1.1.1"

var gGit *gitInfo

// init get remote git
func init() {
    gGit = new(gitInfo)
    gGit.remote = git.NewRemote(
        memory.NewStorage(), &config.RemoteConfig{
            Name: "origin",
            URLs: []string{"https://github.com/go-roc/roc"},
        },
    )
    gGit.enableDynamic = true
    gGit.latestTag = defaultLatestTag
}

// getLatestTag get latest git tag or default tag
func getLatestTag() string {

    if !gGit.enableDynamic {
        return gGit.latestTag
    }

    refs, err := gGit.remote.List(&git.ListOptions{})
    if err != nil {
        rlog.Error(err)
        return gGit.latestTag
    }

    var tags []string
    for _, ref := range refs {
        if ref.Name().IsTag() {
            tags = append(tags, ref.Name().Short())
        }
    }

    if len(tags) > 0 {
        TagsSort(tags)
        gGit.latestTag = tags[len(tags)-1]
    }

    return gGit.latestTag
}

func TagsSort(a []string) {
    sort.Sort(TagSlice(a))
}

type TagSlice []string

func (p TagSlice) Len() int { return len(p) }
func (p TagSlice) Less(i, j int) bool {
    f := func(s string) string {
        return strings.TrimFunc(
            s, func(r rune) bool {
                if r == 'v' || r == '.' {
                    return true
                }
                return false
            },
        )
    }

    return strings.Compare(f(p[i][:]), f(p[j][:])) < 0
}
func (p TagSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
