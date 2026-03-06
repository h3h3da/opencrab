package cli

import (
	"fmt"
	"os"
)

// Command represents a CLI command.
type Command struct {
	Use   string
	Short string
	Long  string
	RunE  func(cmd *Command, args []string) error

	subcommands []*Command
}

// AddCommand adds a subcommand.
func (c *Command) AddCommand(cmd *Command) {
	c.subcommands = append(c.subcommands, cmd)
}

// Execute runs the command or dispatches to appropriate subcommand.
func (c *Command) Execute() error {
	args := os.Args[1:]
	if len(args) > 0 {
		for _, sub := range c.subcommands {
			if sub.Use == args[0] {
				return sub.Execute()
			}
		}
	}
	return c.RunE(c, args)
}

func (c *Command) ExecuteWithArgs(args []string) error {
	if len(args) > 0 {
		for _, sub := range c.subcommands {
			if sub.Use == args[0] {
				return sub.ExecuteWithArgs(args[1:])
			}
		}
	}
	return c.RunE(c, args)
}
