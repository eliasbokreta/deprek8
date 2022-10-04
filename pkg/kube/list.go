package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// nolint: revive
type KubeResource struct {
	Namespace    string `json:"namespace" yaml:"namespace"`
	Name         string `json:"name" yaml:"name"`
	Kind         string `json:"kind" yaml:"kind"`
	GroupVersion string `json:"groupVersion" yaml:"groupVersion"`
}

// List all resources
func List(kubeContext string) (*[]KubeResource, string, error) {
	var deprecatedResources DeprecatedResources
	deprecatedResources.Load()

	tmpKubeconfigPath, err := GenerateTemporaryKubeconfig(kubeContext)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate temporary kubeconfig: %w", err)
	}

	clientConfig, err := CreateConfigClient(tmpKubeconfigPath).ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("could not get clientset: %w", err)
	}

	clientConfig.WarningHandler = rest.NoWarnings{}

	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return nil, "", fmt.Errorf("could not create client config: %w", err)
	}

	gvrs := deprecatedResources.GVRSBuilder()

	krs := []KubeResource{}

	for _, gvr := range gvrs {
		resources, err := dynamicClient.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, resource := range resources.Items {
			if resource.GetNamespace() != "" {
				kr := KubeResource{
					Namespace:    resource.GetNamespace(),
					Name:         resource.GetName(),
					Kind:         resource.GroupVersionKind().Kind,
					GroupVersion: resource.GetObjectKind().GroupVersionKind().GroupVersion().String(),
				}

				if deprecatedResources.GVRSParser(kr) {
					krs = append(krs, kr)
				}
			}
		}
	}

	return &krs, GetServerVersion(clientConfig), nil
}
