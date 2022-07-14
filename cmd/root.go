package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var outputTypes = []string{
	"json",
	"yaml",
	"text",
}

var (
	allNamespaces bool
	outputType    string
	version       string
)

var rootCmd = &cobra.Command{
	Use:     "deprek8",
	Short:   "deprek8 is an utility Kubernetes/Helm tool",
	Long:    "deprek8 is an utility tool that reports deprecated Kubernetes resources and Helm Charts",
	Version: version,
}

var docCmd = &cobra.Command{
	Use:    "doc",
	Short:  "deprek8 cmd documentation",
	Long:   "deprek8 commands' markdown documentation",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := doc.GenMarkdownTree(rootCmd, "./docs"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func initCmd() {
	cobra.OnInitialize()
	deprek8CmdInit()
	rootCmd.AddCommand(docCmd)
}

func Execute() error {
	initCmd()
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("could not run the command tree: %w", err)
	}

	return nil
}
