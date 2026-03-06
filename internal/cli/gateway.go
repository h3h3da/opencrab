package cli

import (
	"fmt"

	"github.com/h3h3da/opencrab/internal/gateway"
)

func newGatewayCommand() *Command {
	return &Command{
		Use:   "gateway",
		Short: "Start the OpenCrab gateway",
		Long: `Start the secure OpenCrab gateway control plane.
By default binds to loopback (127.0.0.1) for security.
Use --port to specify port (default: 18789).`,

		RunE: runGateway,
	}
}

func runGateway(cmd *Command, args []string) error {
	cfg := gateway.DefaultConfig()
	cfg.Port = 18789
	cfg.BindLoopback = true // Security: default to loopback only

	srv, err := gateway.NewServer(cfg)
	if err != nil {
		return fmt.Errorf("create gateway: %w", err)
	}

	fmt.Println("🦀 OpenCrab Gateway starting...")
	fmt.Printf("   Bind: %s (loopback only - secure default)\n", cfg.Address())
	fmt.Printf("   Port: %d\n", cfg.Port)
	fmt.Println("   Auth: required (token or password)")
	fmt.Println("")

	return srv.Run()
}
