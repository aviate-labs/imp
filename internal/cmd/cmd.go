package cmd

import (
	"fmt"
	"strings"
)

type Command struct {
	Name     string
	Aliases  []string
	Summary  string
	Commands []Command
	Method   func(args []string) error
}

func (c Command) Call(args ...string) error {
	if c.Method != nil {
		if err := c.Method(args); err != nil {
			c.Help()
		}
		return nil
	}

	name := args[0]
	if name == "help" {
		c.Help()
		return nil
	}

	var cmd Command
	for _, c := range c.Commands {
		for _, n := range append([]string{c.Name}, c.Aliases...) {
			if n == name {
				cmd = c
				break
			}
		}
	}
	if cmd.Name == "" {
		return fmt.Errorf("unknown")
	}
	return cmd.Call(args[1:]...)
}

func (c Command) Help() {
	fmt.Println(strings.Title(c.Summary))
	fmt.Println()
	fmt.Printf("Usage:\n\t%s", c.Name)
	var cmds []string
	for _, c := range c.Commands {
		cmds = append(cmds, fmt.Sprintf("%s\t\t%s", c.Name, c.Summary))
	}
	if len(cmds) == 0 {
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
