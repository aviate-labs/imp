package cmd

import (
	"fmt"
	"strings"
)

func (c Command) Help() {
	fmt.Println(strings.Title(c.Summary))
	fmt.Println()
	fmt.Printf("Usage:\n\t%s", c.Name)
	var cmds []string
	for _, c := range c.Commands {
		cmds = append(cmds, fmt.Sprintf("%s\t\t%s", c.Name, c.Summary))
	}
	if len(cmds) == 0 {
		var args []string
		for _, a := range c.Args {
			args = append(args, fmt.Sprintf("<%s>", a))
		}
		if len(args) != 0 {
			fmt.Printf(" %s", strings.Join(args, " "))
		}
		fmt.Println()
	} else {
		fmt.Println(" <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Print("\t")
		fmt.Println(strings.Join(cmds, "\n\t"))
	}
	fmt.Println()
}
