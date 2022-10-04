package kube

import (
	"fmt"
	"path/filepath"

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

func GetServerVersion(clientConfig *rest.Config) string {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	if err != nil {
		return ""
	}

	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return ""
	}

	return serverVersion.String()
}

// Create a Configclient
func CreateConfigClient(kubeconfigPath string) clientcmd.ClientConfig {
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(homedir.HomeDir(), kubeconfigDir, kubeconfigFilename)
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{},
	)
}

// Create a Rawconfig
func GetRawAPIConfig(kubeconfigPath string) (*api.Config, error) {
	configclient := CreateConfigClient(kubeconfigPath)

	rawConfig, err := configclient.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get kube rawconfig: %w", err)
	}

	return &rawConfig, nil
}
