package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "go-generate-code"}
	rootCmd.RemoveCommand()
	rootCmd.AddCommand(NewUserCommand())
	rootCmd.Execute()
}
