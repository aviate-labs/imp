package cmd

import (
	"fmt"
	"strings"
)

// A command with either sub-commands or a list of arguments.
type Command struct {
	// The name of the command.
	Name string
	// Aliases of the name.
	// e.g. version -> v
	Aliases []string
	// A summary explaining the function of the command.
	Summary string

	// A list of sub commands.
	Commands []Command

	// A list of arguments.
	Args Arguments
	// Options of the command.
	// e.g. --all, etc.
	Options []Option
	// The method corresponding with a list of arguments.
	Method func(args []string, options map[string]string) error
}

type Arguments []string

type Option struct {
	Name  string
	Value bool
}

func (c Command) Call(args ...string) error {
	if c.Method != nil {
		return c.method(args)
	}

	if len(args) == 0 {
		c.Help()
		return nil
	}
	name := args[0]
	if name == "help" {
		c.Help()
		return nil
	}
	return c.command(name, args[1:])
}

func (c Command) method(args Arguments) error {
	if len(args) == 1 && args[0] == "help" {
		c.Help()
		return nil
	}
	args, opts := c.extractOptions(args)
	if err := c.checkArguments(args); err != nil {
		fmt.Println(err)
		c.Help()
		return nil
	}
	if err := c.Method(args, opts); err != nil {
		return err
	}
	return nil
}

func (c Command) extractOptions(args Arguments) (Arguments, map[string]string) {
	var (
		arguments Arguments
		arg       string
		options   = make(map[string]string)
	)
	for _, a := range args {
		if arg != "" {
			options[arg] = a
			arg = ""
			continue
		}

		if a, ok := trimPrefix(a, "--"); ok {
			var cont bool
			for _, o := range c.Options {
				if a, ok := trimPrefix(a, o.Name); ok {
					if a, ok := trimPrefix(a, "="); ok && o.Value {
						if a != "" {
							options[o.Name] = a
							cont = true
							break
						}
					}
					if a == "" {
						if o.Value {
							arg = o.Name
						} else {
							options[o.Name] = ""
						}
						cont = true
						break
					}
				}
			}
			if cont {
				continue
			}
		}
		arguments = append(arguments, a)
	}
	return arguments, options
}

func trimPrefix(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return strings.TrimPrefix(s, prefix), true
	}
	return s, false
}

func (c Command) command(name string, args Arguments) error {
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
		return fmt.Errorf("command not found")
	}
	if err := cmd.Call(args...); err != nil {
		fmt.Println(err)
		c.Help()
	}
	return nil
}

// checkArguments returns an error if the number of arguments do not equal the
// expected amount.
func (c Command) checkArguments(args Arguments) error {
	l := len(c.Args)
	if len(args) != l {
		var s []string
		for _, a := range c.Args {
			s = append(s, fmt.Sprintf("<%s>", a))
		}

		switch l {
		case 0:
			return fmt.Errorf("expected no argument")
		case 1:
			return fmt.Errorf("expected 1 argument: %s", s[0])
		default:
			return fmt.Errorf("expected %d argument(s): %s", len(c.Args), strings.Join(s, " "))
		}
	}
	return nil
}
