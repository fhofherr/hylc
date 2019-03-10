package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hylc",
	Short: "hylc is a login and consent provider for Ory Hydra",
}

// Execute starts hylc and executes the requested command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v", errors.Wrap(err, "cobra root command"))
		os.Exit(1)
	}
}
