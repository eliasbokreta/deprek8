package deprek8

import (
	"fmt"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/eliasbokreta/deprek8/pkg/helm"
	"github.com/eliasbokreta/deprek8/pkg/kube"
	"github.com/eliasbokreta/deprek8/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Action           string
	OutputType       string
	SelectedContexts []string
	ExportToCSV      bool
}

type FinalOutput struct {
	ClusterName   string
	ServerVersion string
	Resources     interface{}
}

func New(config Config) *Config {
	return &config
}

// Generate a select menu to choose Kubernetes contexts
func (c *Config) SelectContexts() error {
	promptValues := []string{}

	apiconfig, err := kube.GetRawAPIConfig("")
	if err != nil {
		return fmt.Errorf("could not get raw API config: %w", err)
	}

	for cluster := range apiconfig.Contexts {
		promptValues = append(promptValues, cluster)
	}

	prompt := &survey.MultiSelect{
		Message: "Which context(s) do you want to select ?",
		Options: promptValues,
	}

	err = survey.AskOne(prompt, &c.SelectedContexts, survey.WithKeepFilter(true))
	if err != nil {
		return fmt.Errorf("cannot generate survey: %w", err)
	}

	return nil
}

func (c *Config) Run() error {
	if err := c.SelectContexts(); err != nil {
		return fmt.Errorf("could not select new Kube context: %w", err)
	}

	finalOutput := []FinalOutput{}

	switch c.Action {
	case "helm":
		var wg sync.WaitGroup
		for _, kubeContext := range c.SelectedContexts {
			wg.Add(1)
			go func(kctx string) {
				defer wg.Done()

				output := FinalOutput{}
				output.ClusterName = kctx

				h, err := helm.New()
				if err != nil {
					log.Errorf("could not create new Helm object for context '%v': %v", kctx, err)
					return
				}

				output.Resources, output.ServerVersion, err = h.List(kctx)
				if err != nil {
					log.Errorf("could not list Helm releases for context '%v': %v", kctx, err)
					return
				}

				finalOutput = append(finalOutput, output)
			}(kubeContext)
		}
		wg.Wait()

	case "kube":
		var wg sync.WaitGroup
		for _, kubeContext := range c.SelectedContexts {
			wg.Add(1)
			go func(kctx string) {
				defer wg.Done()

				var err error
				output := FinalOutput{}
				output.ClusterName = kctx

				output.Resources, output.ServerVersion, err = kube.List(kctx)
				if err != nil {
					log.Errorf("could not list Kube resources: %v", err)
					return
				}

				finalOutput = append(finalOutput, output)
			}(kubeContext)
		}
		wg.Wait()
	}

	for _, r := range finalOutput {
		log.Infof("Fetching resources from context: '%s' - Cluster version: '%s'", r.ClusterName, r.ServerVersion)
		if err := utils.OutputResult(r.Resources, c.OutputType); err != nil {
			log.Errorf("could not output resources: %v", err)
		}

		if c.ExportToCSV {
			if err := utils.SaveToCSV(r.ClusterName, r.Resources); err != nil {
				log.Errorf("could not save data to file: %v", err)
			}
		}
	}

	return nil
}
