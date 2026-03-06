package cli

import "fmt"

func newVersionCommand() *Command {
	return &Command{
		Use:   "version",
		Short: "Print version information",

		RunE: func(cmd *Command, args []string) error {
			fmt.Println("opencrab v0.1.0")
			fmt.Println("Go rewrite of OpenClaw with enhanced security")
			return nil
		},
	}
}
