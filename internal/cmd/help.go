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

		var options []string
		for _, o := range c.Options {
			s := o.Name
			if o.HasValue {
				s += "\t\t<value>"
			}
			options = append(options, s)
		}
		if len(options) != 0 {
			fmt.Printf("\nArguments:\n\t")
			fmt.Printf("%s", strings.Join(options, "\n\t"))
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
