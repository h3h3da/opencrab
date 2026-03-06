package cli

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Execute runs the CLI application.
func Execute() error {
	if err := setupLogging(); err != nil {
		return err
	}

	rootCmd := newRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("command failed")
		return err
	}
	return nil
}

func setupLogging() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return nil
}

func newRootCommand() *Command {
	cmd := &Command{
		Use:   "opencrab",
		Short: "OpenCrab - Secure personal AI assistant gateway",
		Long: `OpenCrab is a security-hardened Go rewrite of OpenClaw.
Your own personal AI assistant. Any OS. Any Platform. The crab way. 🦀`,

		RunE: func(c *Command, args []string) error {
			fmt.Println("OpenCrab - Secure personal AI assistant")
			fmt.Println("Run 'opencrab gateway' to start the gateway")
			fmt.Println("Run 'opencrab --help' for more options")
			return nil
		},
	}

	cmd.AddCommand(newGatewayCommand())
	cmd.AddCommand(newVersionCommand())
	return cmd
}
