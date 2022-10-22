package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/eliasbokreta/deprek8/pkg/helm"
	"github.com/eliasbokreta/deprek8/pkg/kube"
	log "github.com/sirupsen/logrus"
)

const baseFilename = "deprek8_export"

// Create a temporary directory for saving files
func createTmpDirectory() (string, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "deprek8")
	if err != nil {
		return "", fmt.Errorf("could not create tempfile: %w", err)
	}

	return tmpDir, nil
}

// Dump data to CSV file
func SaveToCSV(filename, data interface{}) error {
	baseDir, err := createTmpDirectory()
	if err != nil {
		return fmt.Errorf("could not create tmpdir: %w", err)
	}

	finalPath := fmt.Sprintf("%s/%s_%s.csv", baseDir, baseFilename, filename)

	csvFile, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	kubeResources, ok := data.(*[]kube.KubeResource)
	if ok {
		headers := []string{"Namespace", "Name", "Kind", "GroupVersion"}
		err := w.Write(headers)
		if err != nil {
			return fmt.Errorf("could not write headers to csv file: %w", err)
		}

		for _, resource := range *kubeResources {
			line := []string{
				resource.Namespace,
				resource.Name,
				resource.Kind,
				resource.GroupVersion,
			}

			err := w.Write(line)
			if err != nil {
				return fmt.Errorf("could not write data to csv file: %w", err)
			}
		}
		log.Infof("File successfully saved to '%s'", finalPath)
		return nil
	}

	helmReleases, ok := data.(*[]helm.Release)
	if ok {
		headers := []string{"Namespace", "Repository", "Name", "Version", "Latest", "Deprecated Resources"}
		err := w.Write(headers)
		if err != nil {
			return fmt.Errorf("could not write headers to csv file: %w", err)
		}

		for _, release := range *helmReleases {
			line := []string{
				release.Namespace,
				release.Repository,
				release.ChartName,
				release.ChartVersion,
				release.LatestChartVersion,
				release.DeprecatedObjects,
			}

			err := w.Write(line)
			if err != nil {
				return fmt.Errorf("could not write data to csv file: %w", err)
			}
		}
		log.Infof("File successfully saved to '%s'", finalPath)
		return nil
	}

	return nil
}
