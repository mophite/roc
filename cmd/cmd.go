package cmd

import (
    "fmt"
    "os"
    "sort"

    "github.com/urfave/cli/v2"
)

var gCmd *cmd

type cmd struct {
    *cli.App
}

func Register(c ...*cli.Command) {
    if gCmd == nil {
        gCmd = new(cmd)
        gCmd.App = cli.NewApp()
    }

    gCmd.Commands = append(gCmd.Commands, c...)

    sort.Sort(cli.CommandsByName(gCmd.Commands))
}

func Run() {
    err := gCmd.Run(os.Args)
    if err != nil {
        fmt.Fprintln(os.Stderr, "roc cmd err ", err.Error())
    }
}
