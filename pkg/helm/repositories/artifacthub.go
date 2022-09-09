package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

const (
	searchEndpoint = "/packages/search"
)

type ArtifactHub struct {
	BaseURL string
}

type SearchPackageResponse struct {
	Packages []struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Repository struct {
			Name string `json:"name"`
		} `json:"repository"`
	} `json:"packages"`
}

func NewArtifacthub() *ArtifactHub {
	return &ArtifactHub{
		BaseURL: viper.GetString("deprek8.helm.repository.artifacthub.apiBaseUrl"),
	}
}

func (a *ArtifactHub) requestArtifacthub(chartName string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	urlValues := url.Values{}
	urlValues.Add("ts_query_web", chartName)
	urlValues.Add("limit", "1")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", a.BaseURL, searchEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("could not prepare http request: %w", err)
	}

	q := req.URL.Query()
	for key, values := range urlValues {
		for _, value := range values {
			q.Add(key, value)
		}
	}

	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while executing request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading request body: %w", err)
	}

	return body, nil
}

func (a *ArtifactHub) GetChartLatestVersion(chartName string) (string, error) {
	body, err := a.requestArtifacthub(chartName)
	if err != nil {
		return "", fmt.Errorf("error while trying to request artifacthub: %w", err)
	}

	searchPackageResponse := SearchPackageResponse{}
	if err := json.Unmarshal(body, &searchPackageResponse); err != nil {
		return "", fmt.Errorf("could not unmarshal search package response: %w", err)
	}

	if len(searchPackageResponse.Packages) == 0 {
		return "", nil
	}

	return searchPackageResponse.Packages[0].Version, nil
}
