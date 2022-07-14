package cmd

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/deprek8/pkg/deprek8"
	"github.com/eliasbokreta/deprek8/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	filterDeprecated bool
	filterChartName  string
)

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "list helm releases",
	Long:  `list helm releases in the current Kubernetes context`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !utils.StringInSlice(outputType, outputTypes) {
			return fmt.Errorf("--output must be one of %v", outputTypes)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := deprek8.Config{
			Action:           "helm",
			AllNamespaces:    allNamespaces,
			OutputType:       outputType,
			FilterDeprecated: filterDeprecated,
			FilterChartName:  filterChartName,
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
	Short: "list Kube resources",
	Long:  "list deprecated Kubernetes resources",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := deprek8.Config{
			Action:     "kube",
			OutputType: outputType,
		}
		d := deprek8.New(cfg)

		if err := d.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "select Kube context",
	Long:  "select a Kubernetes context",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := deprek8.Config{
			Action: "context",
		}
		d := deprek8.New(cfg)

		if err := d.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func deprek8CmdInit() {
	rootCmd.AddCommand(helmCmd)

	kubeCmd.Flags().StringVarP(&outputType, "output", "o", "text", "Choose type of output (json|yaml|text)")

	helmCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", false, "Fetch data on all namespaces")
	helmCmd.Flags().BoolVarP(&filterDeprecated, "filter-deprecated", "d", false, "Filter helm releases with deprecated k8s resources")
	helmCmd.Flags().StringVarP(&filterChartName, "filter-name", "n", "", "Filter Helm chart's name")

	rootCmd.AddCommand(kubeCmd)
	rootCmd.AddCommand(contextCmd)
}
