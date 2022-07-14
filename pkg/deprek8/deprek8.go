package deprek8

import (
	"fmt"

	"github.com/eliasbokreta/deprek8/pkg/helm"
	"github.com/eliasbokreta/deprek8/pkg/kube"
	"github.com/eliasbokreta/deprek8/pkg/utils"
)

type Config struct {
	Action           string
	OutputType       string
	AllNamespaces    bool
	FilterDeprecated bool
	FilterChartName  string
}

func New(config Config) *Config {
	return &config
}

// nolint: cyclop
func (c *Config) Run() error {
	switch c.Action {
	case "helm":
		h, err := helm.New()
		if err != nil {
			return fmt.Errorf("could not create new Helm object: %w", err)
		}

		releases, err := h.List(c.AllNamespaces)
		if err != nil {
			return fmt.Errorf("could not list helm releases: %w", err)
		}

		if c.FilterDeprecated {
			releases = utils.FilterDeprecatedReleases(*releases)
		}

		if c.FilterChartName != "" {
			releases = utils.FilterPerChartName(*releases, c.FilterChartName)
		}

		if err := utils.OutputResult(releases, c.OutputType); err != nil {
			return fmt.Errorf("could not output helm releases: %w", err)
		}
	case "kube":
		k, err := kube.New()
		if err != nil {
			return fmt.Errorf("could not create new Kube object: %w", err)
		}

		out, err := k.List()
		if err != nil {
			return fmt.Errorf("could not list Kube resources: %w", err)
		}

		if err := utils.OutputResult(out, c.OutputType); err != nil {
			return fmt.Errorf("could not output kube resources: %w", err)
		}
	case "context":
		k, err := kube.New()
		if err != nil {
			return fmt.Errorf("could not create new Kube object: %w", err)
		}

		if err := k.SelectContext(); err != nil {
			return fmt.Errorf("could not select new Kube context: %w", err)
		}
	}

	return nil
}
