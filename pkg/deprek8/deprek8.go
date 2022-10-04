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
	AllNamespaces    bool
	FilterDeprecated bool
	FilterChartName  string
	SelectedContexts []string
}

type FinalOutput struct {
	ClusterName   string
	ServerVersion string
	Releases      *[]helm.Release
	KubeResources *[]kube.KubeResource
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

// nolint: cyclop,funlen,gocognit
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

				h, err := helm.New()
				if err != nil {
					log.Errorf("could not create new Helm object for context '%v': %v", kctx, err)
					return
				}

				releases, serverVersion, err := h.List(kctx, c.AllNamespaces)
				if err != nil {
					log.Errorf("could not list helm releases for context '%v': %v", kctx, err)
					return
				}

				if c.FilterDeprecated {
					releases = utils.FilterDeprecatedReleases(*releases)
				}

				if c.FilterChartName != "" {
					releases = utils.FilterPerChartName(*releases, c.FilterChartName)
				}

				output.ClusterName = kctx
				output.ServerVersion = serverVersion
				output.Releases = releases

				finalOutput = append(finalOutput, output)
			}(kubeContext)
		}
		wg.Wait()

		for _, r := range finalOutput {
			log.Infof("Fetching helm releases from context: '%s' - Cluster version: '%s'", r.ClusterName, r.ServerVersion)
			if err := utils.OutputResult(r.Releases, c.OutputType); err != nil {
				log.Errorf("could not output helm releases: %v", err)
			}
		}
	case "kube":
		var wg sync.WaitGroup
		for _, kubeContext := range c.SelectedContexts {
			wg.Add(1)
			go func(kctx string) {
				defer wg.Done()

				output := FinalOutput{}

				resources, serverVersion, err := kube.List(kctx)
				if err != nil {
					log.Errorf("could not list Kube resources: %v", err)
					return
				}

				output.ClusterName = kctx
				output.ServerVersion = serverVersion
				output.KubeResources = resources

				finalOutput = append(finalOutput, output)
			}(kubeContext)
		}
		wg.Wait()

		for _, r := range finalOutput {
			log.Infof("Fetching Kubernetes resources from context: '%s' - Cluster version: '%s'", r.ClusterName, r.ServerVersion)
			if err := utils.OutputResult(r.KubeResources, c.OutputType); err != nil {
				return fmt.Errorf("could not output kube resources: %w", err)
			}
		}
	}

	return nil
}
