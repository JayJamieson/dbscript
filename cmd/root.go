package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbscript",
	Short: "Easily script your CDC events",
	Long: `Script your CDC events with Javascript. Build complex or simple pre-processing
pipelines in Javascript.

Javscript files can be hot reloading to avoid downtime.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {

	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
