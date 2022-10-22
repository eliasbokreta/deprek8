package cmd

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/deprek8/pkg/deprek8"
	"github.com/eliasbokreta/deprek8/pkg/utils"
	"github.com/spf13/cobra"
)

var outputTypes = []string{
	"json",
	"yaml",
	"text",
}

var (
	outputType string
	export     bool
)

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "list helm releases",
	Long:  `list helm releases with deprecated objects`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !utils.StringInSlice(outputType, outputTypes) {
			return fmt.Errorf("--output must be one of %v", outputTypes)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := deprek8.Config{
			Action:      "helm",
			OutputType:  outputType,
			ExportToCSV: export,
		}
		d := deprek8.New(cfg)

		if err := d.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "list Kube objects",
	Long:  "list deprecated Kubernetes objects",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := deprek8.Config{
			Action:      "kube",
			OutputType:  outputType,
			ExportToCSV: export,
		}
		d := deprek8.New(cfg)

		if err := d.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func deprek8CmdInit() {
	helmCmd.Flags().StringVarP(&outputType, "output", "o", "text", "Choose type of output (json|yaml|text)")
	helmCmd.Flags().BoolVarP(&export, "export", "e", false, "Save output to csv file")
	rootCmd.AddCommand(helmCmd)

	kubeCmd.Flags().StringVarP(&outputType, "output", "o", "text", "Choose type of output (json|yaml|text)")
	kubeCmd.Flags().BoolVarP(&export, "export", "e", false, "Save output to csv file")
	rootCmd.AddCommand(kubeCmd)
}
