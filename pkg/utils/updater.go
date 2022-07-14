package utils

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	log "github.com/sirupsen/logrus"
)

// Update the binary version if a newer is found
func SelfUpdate(version string) error {
	version = version[1:]

	latest, found, err := selfupdate.DetectLatest("eliasbokreta/deprek8")
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}

	v := semver.MustParse(version) // nolint: ifshort
	if !found || latest.Version.LTE(v) {
		log.Info("Current version is the latest")
		return nil
	}

	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Do you like to update to version '%s' ?", latest.Version),
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return fmt.Errorf("could not generate survey: %w", err)
	}

	if !confirm {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not locate executable path: %w", err)
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	log.Infof("Successfully updated to version '%v'", latest.Version)

	return nil
}
