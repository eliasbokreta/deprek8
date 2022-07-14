package kube

import (
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth/azure" // Import needed for Azure authentication
)

const (
	kubeconfigDir      = ".kube"
	kubeconfigFilename = "config"
)

type Kube struct {
	Config          api.Config
	ClientConfig    *rest.Config
	Contexts        []string
	SelectedContext string
	ClusterVersion  string
}

// Create a new Kube struct containing configuration and contexts
func New() (*Kube, error) {
	kubeConfigPath := filepath.Join(homedir.HomeDir(), kubeconfigDir, kubeconfigFilename)
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{},
	)

	clientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get clientconfig: %w", err)
	}

	rawConfig, err := config.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get kube rawconfig: %w", err)
	}

	contexts := []string{}
	for cluster := range rawConfig.Contexts {
		contexts = append(contexts, cluster)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create discovery client: %w", err)
	}

	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("could not get server version: %w", err)
	}

	return &Kube{
		Config:          rawConfig,
		ClientConfig:    clientConfig,
		Contexts:        contexts,
		SelectedContext: rawConfig.CurrentContext,
		ClusterVersion:  fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor),
	}, nil
}

// Select a context in an interactive manner
func (k *Kube) SelectContext() error {
	promptValues := k.Contexts

	var outputValue string

	prompt := &survey.Select{
		Message: "Choose your Kubernetes context",
		Options: promptValues,
	}

	err := survey.AskOne(prompt, &outputValue, survey.WithKeepFilter(true))
	if err != nil {
		return fmt.Errorf("unable to prompt survey: %w", err)
	}

	k.SelectedContext = outputValue
	if err := k.switchContext(); err != nil {
		return fmt.Errorf("unable to switch context: %w", err)
	}

	return nil
}

// Switch the active context to the selected one
func (k *Kube) switchContext() error {
	if k.Config.Contexts[k.SelectedContext] == nil {
		return fmt.Errorf("context %s doesn't exists", k.SelectedContext)
	}
	k.Config.CurrentContext = k.SelectedContext
	err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), k.Config, true)
	if err != nil {
		return fmt.Errorf("could not modify kubeconfig: %w", err)
	}

	log.Infof("Successfully switched to '%s' context", k.SelectedContext)

	return nil
}
