package helm

import (
	"fmt"

	"github.com/eliasbokreta/deprek8/pkg/helm/repositories"
	"github.com/eliasbokreta/deprek8/pkg/kube"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

type Helm struct {
	Releases    []Release `json:"releases" yaml:"releases"`
	Kube        *kube.Kube
	ArtifactHub *repositories.ArtifactHub
}

type Release struct {
	Namespace          string           `json:"namespace" yaml:"namespace"`
	ReleaseName        string           `json:"releaseName" yaml:"releaseName"`
	ChartName          string           `json:"chartName" yaml:"chartName"`
	ChartVersion       string           `json:"chartVersion" yaml:"chartVersion"`
	LatestChartVersion string           `json:"latestChartVersion" yaml:"latestChartVersion"`
	Repository         string           `json:"repository" yaml:"repository"`
	ChartRevision      int              `json:"chartRevision" yaml:"chartRevision"`
	ChartStatus        string           `json:"chartStatus" yaml:"chartStatus"`
	Resources          []*kube.Resource `json:"deprecatedResources,omitempty" yaml:"deprecatedResources,omitempty"`
}

func New() (*Helm, error) {
	k, err := kube.New()
	if err != nil {
		return nil, fmt.Errorf("could not instantiate Kube: %w", err)
	}

	return &Helm{
		Releases:    []Release{},
		Kube:        k,
		ArtifactHub: repositories.NewArtifacthub(),
	}, nil
}

// nolint: ineffassign, wastedassign
func (h *Helm) getLatestChartVersion(rel *release.Release, release *Release) error {
	latestVersion := "unknown"
	repository := "unknown"

	latestVersion, err := h.ArtifactHub.GetChartLatestVersion(rel.Chart.Name())
	if err != nil {
		return fmt.Errorf("could not fetch latest chart version from Artifacthub: %w", err)
	}

	if latestVersion != "" {
		repository = "ArtifactHub"
	}

	release.LatestChartVersion = latestVersion
	release.Repository = repository

	return nil
}

// List all the deployed helm charts in the current namespace or all namespaces in the current context
func (h *Helm) List(allNamespaces bool) (*[]Release, error) {
	log.Infof("Fetching helm releases from context: '%s'", h.Kube.Config.CurrentContext)

	var deprecatedResources kube.DeprecatedResources
	deprecatedResources.Load()

	settings := cli.New()
	actionConfig := new(action.Configuration)
	namespace := ""
	if !allNamespaces {
		namespace = settings.Namespace()
	}

	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "", log.Infof); err != nil {
		return nil, fmt.Errorf("could not init helm action configuration: %w", err)
	}

	client := action.NewList(actionConfig)
	client.Deployed = true

	results, err := client.Run()
	if err != nil {
		return nil, fmt.Errorf("could not run helm list command: %w", err)
	}

	for _, rel := range results {
		release := Release{
			Namespace:     rel.Namespace,
			ReleaseName:   rel.Name,
			ChartName:     rel.Chart.Name(),
			ChartVersion:  rel.Chart.Metadata.Version,
			ChartRevision: rel.Version,
			ChartStatus:   rel.Info.Status.String(),
		}

		// Fetch latest chart version from repositories
		if err := h.getLatestChartVersion(rel, &release); err != nil {
			return nil, fmt.Errorf("could not fetch latest chart version: %w", err)
		}

		// Parse deprecated resources from found manifests
		resources := deprecatedResources.ManifestParser(rel.Manifest)

		release.Resources = resources
		h.Releases = append(h.Releases, release)
	}

	return &h.Releases, nil
}
