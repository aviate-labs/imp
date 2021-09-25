package main

import (
	"fmt"
	"os"

	"github.com/aviate-labs/imp/internal/cmd"
)

var pwd string

func init() {
	pwd, _ = os.Getwd()
}

var imp = cmd.Command{
	Name:    "imp",
	Summary: "experimental command line tool for the Internet Computer",
	Commands: []cmd.Command{
		version,
		stats,
	},
}

var version = cmd.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Summary: "print Imp version",
	Method: func(args []string, _ map[string]string) error {
		if len(args) != 0 {
			return fmt.Errorf("too long")
		}
		fmt.Println("v0.1.0")
		return nil
	},
}

func main() {
	if len(os.Args) == 1 {
		imp.Help()
		return
	}
	if err := imp.Call(os.Args[1:]...); err != nil {
		panic(err)
	}
}
