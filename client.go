package govdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TODO:
// - Add support of context.Context to each of method.

// DefaultURL is a URL of official Ukrainian government data platform.
const DefaultURL = "https://data.gov.ua"

// BaseURL is current base URL used for sending requests.
var BaseURL = DefaultURL

// Client is wrapper of http.Client for Ukrainian government data platform.
type Client struct {
	http.Client
}

// DefaultClient is a default HTTP client.
var DefaultClient = NewClient()

// NewClient creates new instance of client with timeout equal 5 seconds.
func NewClient() *Client {
	return &Client{
		Client: http.Client{
			Timeout: time.Minute * 5,
		},
	}
}

// PackageShow returns information about package by it's unique id.
// For more information: https://data.gov.ua/pages/aboutuser2.
func PackageShow(id string) (*Package, error) {
	return DefaultClient.PackageShow(context.Background(), id)
}

// ResourceShow returns information about resource by it's unique id.
// For more information: https://data.gov.ua/pages/aboutuser2.
func ResourceShow(id string) (*Resource, error) {
	return DefaultClient.ResourceShow(context.Background(), id)
}

// ResourceRevision returns information about specific revision of a specific resource.
// This is not a part of API, but it pretty important method
// for downloading updated version of a resource.
func ResourceRevision(pkg, resource, revision string) (io.ReadCloser, error) {
	return DefaultClient.ResourceRevision(context.Background(), pkg, resource, revision)
}

// PackageShow returns information about package by it's unique id.
// For more information: https://data.gov.ua/pages/aboutuser2.
func (client *Client) PackageShow(ctx context.Context, id string) (*Package, error) {
	url := BaseURL + "/api/3/action/package_show?id=" + id

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var pkg Package
	if err := json.Unmarshal(response.Result, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package: %w", err)
	}

	return &pkg, nil
}

// ResourceShow returns information about resource by it's unique id.
// For more information: https://data.gov.ua/pages/aboutuser2.
func (client *Client) ResourceShow(ctx context.Context, id string) (*Resource, error) {
	url := BaseURL + "/api/3/action/resource_show?id=" + id

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var resource Resource
	if err := json.Unmarshal(response.Result, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	for i := range resource.Revisions {
		parts := strings.Split(resource.Revisions[i].URL, "/")
		resource.Revisions[i].ID = parts[len(parts)-1]
	}

	return &resource, nil
}

// ResourceRevision returns information about specific revision of a specific resource.
// This is not a part of API, but it pretty important method
// for downloading updated version of a resource.
func (client *Client) ResourceRevision(ctx context.Context, pkg, resource, revision string) (io.ReadCloser, error) {
	url := BaseURL + "/dataset/" + pkg + "/resource/" + resource + "/revision/" + revision

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp.Body, nil
}
