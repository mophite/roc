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

package rocx

import (
    "fmt"
    "os"
    "runtime"

    cmdx "github.com/go-cmd/cmd"
    "github.com/go-roc/roc/cmd"
    "github.com/urfave/cli/v2"
)

//init project
func init() {
    cmd.Register(
        &cli.Command{
            Name:        "init",
            Aliases:     []string{"i"},
            Usage:       "Init your project at root dir",
            UsageText:   "",
            Description: "",
            Action:      run,
        },
    )
}

type rocx struct {

    //root dir folder name
    root string

    //to create moduleName
    moduleName string

    //absolute path here
    absolutePath string

    //create file by tpl
    stencil []stencil

    //create folder path by path
    absPath []absPath

    //roc version at latest tag
    rocVersion string
}

type stencil struct {
    path string
    tpl  string
}

type absPath struct {
    path string
}

func run(c *cli.Context) error {
    moduleName := c.Args().First()

    if len(c.Args().First()) == 0 {
        moduleName = getProjectName()
        return nil
    }

    r := &rocx{
        root:         getProjectName(),
        moduleName:   moduleName,
        absolutePath: getLastPwd(),
        rocVersion:   getLatestTag(),
    }

    if os.Getenv("GO111MODULE") == "off" {
        err := os.Setenv("GO111MODULE", "on")
        if err != nil {
            fmt.Fprintln(os.Stderr, "env GO111MODULE incorrect")
            fmt.Fprintln(os.Stdout, "set GO111MODULE")
            // Create Cmd, buffered output
            envCmd := cmdx.NewCmd("go", "env", "-w", "GO111MODULE", "=", "on")

            select {
            case <-envCmd.Start():
            }
        }
    }

    fmt.Fprintln(os.Stdout, "suggest go version is newer than what 1.16.")
    fmt.Fprintln(os.Stdout, "your go version is: ", runtime.Version())

    r.absPath = []absPath{

    }

    r.stencil = []stencil{

    }

    return nil
}

func (r *rocx) createDir() {

}

func (r *rocx) createFile() {

}
