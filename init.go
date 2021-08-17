package main

import (
	"fmt"
	"os"

	"github.com/aviate-labs/imp/internal/cmd"
)

var initialize = cmd.Command{
	Name:    "init",
	Summary: "initialize new module in current directory",
	Args:    []string{"module-name"},
	Method: func(args []string, _ map[string]string) error {
		if _, err := os.Stat(pwd + "/mo.mod"); err == nil {
			fmt.Println("mo.mod already exists")
			return nil
		}
		f, err := os.Create(pwd + "/mo.mod")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("module %s\n\n", args[0]))
		return nil
	},
}
