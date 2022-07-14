package utils

import (
	"strings"

	"github.com/eliasbokreta/deprek8/pkg/helm"
)

// Returns only deprecated releases
func FilterDeprecatedReleases(r []helm.Release) *[]helm.Release {
	out := make([]helm.Release, 0)

	for _, release := range r {
		for _, resource := range release.Resources {
			if resource != nil {
				out = append(out, release)
			}
		}
	}

	return &out
}

// Returns only specific Charts
func FilterPerChartName(r []helm.Release, chartName string) *[]helm.Release {
	out := make([]helm.Release, 0)

	for _, release := range r {
		if strings.EqualFold(release.ChartName, chartName) {
			out = append(out, release)
		}
	}

	return &out
}
