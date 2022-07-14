package utils

import (
	"testing"

	"github.com/eliasbokreta/deprek8/pkg/helm"
	"github.com/eliasbokreta/deprek8/pkg/kube"
)

func TestFilterDeprecatedReleases(t *testing.T) {
	releases := []helm.Release{}

	release := helm.Release{
		ReleaseName:        "MyDeprecatedRelease",
		ChartName:          "MyDeprecatedChart",
		ChartVersion:       "1.0.0",
		LatestChartVersion: "1.1.1",
		Repository:         "Artifacthub",
		ChartRevision:      1,
		ChartStatus:        "Deployed",
	}

	resource := new(kube.Resource)

	resource.Kind = "Ingress"
	resource.APIVersion = "v1beta"
	resource.Metadata.Name = "MyName"
	resource.DeprecatedResource = kube.DeprecatedResource{
		GroupVersion:       "extensions/v1beta1",
		NewGroupVersion:    "networking.k8s.io/v1",
		DeprecationVersion: "1.14",
		RemovalVersion:     "1.22",
		BreakingChange:     true,
		Details: `
- spec.backend is renamed to spec.defaultBackend
- The backend serviceName field is renamed to service.name
- Numeric backend servicePort fields are renamed to service.port.number
- String backend servicePort fields are renamed to service.port.name
- pathType is now required for each specified path. Options are Prefix, Exact, and ImplementationSpecific. To match the undefined v1beta1 behavior, use ImplementationSpecific.`,
	}
	resource.Error = ""
	release.Resources = append(release.Resources, resource)

	releases = append(releases, release)

	release = helm.Release{
		ReleaseName:        "MyRelease",
		ChartName:          "MyChart",
		ChartVersion:       "2.0.0",
		LatestChartVersion: "2.1.1",
		Repository:         "Artifacthub",
		ChartRevision:      5,
		ChartStatus:        "Deployed",
	}

	releases = append(releases, release)

	if len(releases) != 2 {
		t.Fatal("Should have 2 releases in the slice")
	}

	releases = *FilterDeprecatedReleases(releases)

	if len(releases) != 1 {
		t.Fatal("Should have 1 release in the slice")
	}
}

func TestFilterPerChartName(t *testing.T) {
	releases := []helm.Release{}

	release := helm.Release{
		ReleaseName:        "MyDeprecatedRelease",
		ChartName:          "MyDeprecatedChart",
		ChartVersion:       "1.0.0",
		LatestChartVersion: "1.1.1",
		Repository:         "Artifacthub",
		ChartRevision:      1,
		ChartStatus:        "Deployed",
	}

	releases = append(releases, release)

	release = helm.Release{
		ReleaseName:        "MyRelease",
		ChartName:          "MyChart",
		ChartVersion:       "2.0.0",
		LatestChartVersion: "2.1.1",
		Repository:         "MyRepository",
		ChartRevision:      5,
		ChartStatus:        "Deployed",
	}

	releases = append(releases, release)

	if len(releases) != 2 {
		t.Fatalf("Should have 2 releases in the slice")
	}

	releases = *FilterPerChartName(releases, "MyDeprecatedChart")

	if len(releases) != 1 {
		t.Fatalf("Should have 1 release in the slice")
	}
}
