package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KenichiTanino/archive-tar-go/archive"
)

var rootCmd = &cobra.Command{
	Use: "archive-tar-go",
	Run: func(cmd *cobra.Command, args []string) {
		input_dir := args[0]
		output_file := args[1]
		err := archive.Create(input_dir, output_file)
		if err != nil {
			fmt.Printf("err %s", err)
		}
	},
}

// Execute start command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// print.Fatal(err)
	}

}
