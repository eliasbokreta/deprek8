package utils

import (
	"encoding/json"
	"fmt"

	"github.com/eliasbokreta/deprek8/pkg/helm"
	"github.com/eliasbokreta/deprek8/pkg/kube"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

const (
	red = "red"
)

// Output any struct to json or yaml
func OutputResult(data interface{}, outputType string) error {
	var output []byte
	var err error

	switch outputType {
	case "json":
		output, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("could not unmarshal json output: %w", err)
		}
	case "yaml":
		output, err = yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("could not unmarshal json output: %w", err)
		}
	case "text":
		releases, ok := data.(*[]helm.Release)
		if ok {
			OutputHelm(*releases)
			return nil
		}

		kube, ok := data.(*[]kube.KubeResource)
		if ok {
			OutputKube(*kube)
			return nil
		}

		return nil
	default:
		return fmt.Errorf("unknown outputType value: %s", outputType)
	}

	fmt.Println(string(output))

	return nil
}

func OutputHelm(releases []helm.Release) {
	table := GetTableWriter([]string{"Namespace", "Repository", "Name", "Version", "Latest", "Deprecated Resources"})

	for _, release := range releases {
		line := []string{
			release.Namespace,
			release.Repository,
			release.ChartName,
			release.ChartVersion,
			release.LatestChartVersion,
			release.DeprecatedObjects,
		}

		colors := []tablewriter.Colors{
			GetTableRowColor(red),
			GetTableRowColor(red),
			GetTableRowColor(red),
			GetTableRowColor(red),
			GetTableRowColor(red),
			GetTableRowColor(red),
		}

		table.Rich(line, colors)
	}
	table.Render()
}

func OutputKube(krs []kube.KubeResource) {
	table := GetTableWriter([]string{"Namespace", "Name", "Kind", "GroupVersion"})
	for _, resource := range krs {
		line := []string{
			resource.Namespace,
			resource.Name,
			resource.Kind,
			resource.GroupVersion,
		}

		colors := []tablewriter.Colors{
			GetTableRowColor(red),
			GetTableRowColor(red),
			GetTableRowColor(red),
			GetTableRowColor(red),
		}

		table.Rich(line, colors)
	}
	table.Render()
}
