package kube

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type DeprecatedResources struct {
	DeprecatedResources map[string][]DeprecatedResource `json:"resources" yaml:"resources"`
}

type DeprecatedResource struct {
	GroupVersion       string `json:"groupVersion" yaml:"groupVersion"`
	APIResource        string `json:"apiResource" yaml:"apiResource"`
	NewGroupVersion    string `json:"newGroupVersion" yaml:"newGroupVersion"`
	DeprecationVersion string `json:"deprecationVersion" yaml:"deprecationVersion"`
	RemovalVersion     string `json:"removalVersion" yaml:"removalVersion"`
	BreakingChange     bool   `json:"breakingChange" yaml:"breakingChange"`
	Details            string `json:"details,omitempty" yaml:"details,omitempty"`
}

type Resource struct {
	Kind       string `json:"kind" yaml:"kind"`
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Metadata   struct {
		Name string `json:"name" yaml:"name"`
	} `json:"metadata" yaml:"metadata"`
	DeprecatedResource DeprecatedResource `json:"deprecatedResource,omitempty" yaml:"deprecatedResource,omitempty"`
	Error              string             `json:"error,omitempty" yaml:"error,omitempty"`
}

// Parse a release to get its Kubernetes resources and verify if any are deprecated
func (d *DeprecatedResources) ManifestParser(m string) []*Resource {
	manifests := strings.Split(m, "---\n")

	var resources []*Resource
	for _, manifest := range manifests {
		if manifest != "" {
			resource := new(Resource)

			// Handles cases where generated helm chart yaml is not valid
			if err := yaml.Unmarshal([]byte(manifest), &resource); err != nil {
				resource.Error = err.Error()
			}

			if resource.Kind != "" {
				for _, rs := range d.DeprecatedResources[resource.Kind] {
					if rs.GroupVersion == resource.APIVersion {
						resource.DeprecatedResource = rs
						resources = append(resources, resource)
					}
				}
			}
		}
	}

	return resources
}

func (d *DeprecatedResources) GVRSBuilder() []schema.GroupVersionResource {
	gvrs := []schema.GroupVersionResource{}

	for _, rs := range d.DeprecatedResources {
		for _, r := range rs {
			gv := strings.Split(r.GroupVersion, "/")
			gvrs = append(gvrs, schema.GroupVersionResource{
				Group:    gv[0],
				Version:  gv[1],
				Resource: r.APIResource,
			})
		}
	}

	return gvrs
}

// Parse GVRS
func (d *DeprecatedResources) GVRSParser(r KubeResource) bool {
	for _, rs := range d.DeprecatedResources[r.Kind] {
		if r.GroupVersion == rs.GroupVersion {
			return true
		}
	}

	return false
}

// Load the deprecated API resources from YAML
func (d *DeprecatedResources) Load() {
	file, err := os.ReadFile(viper.GetString("deprek8.kube.resourcesFile"))
	if err != nil {
		log.Errorf("Could not open file: %v", err)
		return
	}

	if err := yaml.Unmarshal(file, &d); err != nil {
		log.Errorf("Could not unmarshal file: %v", err)
		return
	}
}
