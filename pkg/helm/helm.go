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
	ArtifactHub *repositories.ArtifactHub
}

type Release struct {
	Namespace           string           `json:"namespace" yaml:"namespace"`
	ReleaseName         string           `json:"releaseName" yaml:"releaseName"`
	ChartName           string           `json:"chartName" yaml:"chartName"`
	ChartVersion        string           `json:"chartVersion" yaml:"chartVersion"`
	LatestChartVersion  string           `json:"latestChartVersion" yaml:"latestChartVersion"`
	Repository          string           `json:"repository" yaml:"repository"`
	ChartRevision       int              `json:"chartRevision" yaml:"chartRevision"`
	ChartStatus         string           `json:"chartStatus" yaml:"chartStatus"`
	DeprecatedResources []*kube.Resource `json:"deprecatedResources,omitempty" yaml:"deprecatedResources,omitempty"`
	DeprecatedObjects   string           `json:"deprecatedObjects,omitempty" yaml:"deprecatedObjects,omitempty"`
}

func New() (*Helm, error) {
	return &Helm{
		Releases:    []Release{},
		ArtifactHub: repositories.NewArtifacthub(),
	}, nil
}

// nolint: ineffassign, wastedassign
func (h *Helm) getLatestChartVersion(rel *release.Release, release *Release) error {
	repository := "ArtifactHub"
	latestVersion, err := h.ArtifactHub.GetChartLatestVersion(rel.Chart.Name())
	if err != nil {
		return fmt.Errorf("could not fetch latest chart version from Artifacthub: %w", err)
	}

	if latestVersion == "" {
		repository = "unknown"
		latestVersion = "unknown"
	}

	release.LatestChartVersion = latestVersion
	release.Repository = repository

	return nil
}

// List all the deployed helm charts in the current namespace or all namespaces in the current context
func (h *Helm) List(kubeContext string) (*[]Release, string, error) {
	var deprecatedResources kube.DeprecatedResources
	deprecatedResources.Load()

	tmpKubeconfigPath, err := kube.GenerateTemporaryKubeconfig(kubeContext)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate temporary kubeconfig: %w", err)
	}

	clientConfig, err := kube.CreateConfigClient(tmpKubeconfigPath).ClientConfig()
	if err != nil {
		log.Errorf("could not get clientset: %v", err)
	}

	settings := cli.New()
	settings.KubeConfig = tmpKubeconfigPath
	settings.KubeContext = kubeContext
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), "", "", log.Infof); err != nil {
		return nil, "", fmt.Errorf("could not init helm action configuration: %w", err)
	}

	client := action.NewList(actionConfig)
	client.Deployed = true

	results, err := client.Run()
	if err != nil {
		return nil, "", fmt.Errorf("could not run helm list command: %w", err)
	}

	for _, rel := range results {
		// Parse deprecated resources from found manifests
		deprecatedResources := deprecatedResources.ManifestParser(rel.Manifest)

		if len(deprecatedResources) > 0 {
			release := Release{
				Namespace:     rel.Namespace,
				ReleaseName:   rel.Name,
				ChartName:     rel.Chart.Name(),
				ChartVersion:  rel.Chart.Metadata.Version,
				ChartRevision: rel.Version,
				ChartStatus:   rel.Info.Status.String(),
			}

			for i, deprecatedResource := range deprecatedResources {
				release.DeprecatedObjects += fmt.Sprintf("%s(%s)", deprecatedResource.Kind, deprecatedResource.DeprecatedResource.RemovalVersion)
				if i < (len(deprecatedResources) - 1) {
					release.DeprecatedObjects += " "
				}
			}

			// Fetch latest chart version from repositories
			if err := h.getLatestChartVersion(rel, &release); err != nil {
				return nil, "", fmt.Errorf("could not fetch latest chart version: %w", err)
			}

			release.DeprecatedResources = deprecatedResources
			h.Releases = append(h.Releases, release)
		}
	}

	return &h.Releases, kube.GetServerVersion(clientConfig), nil
}
