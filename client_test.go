package govdata

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Fixture struct {
	t     *testing.T
	Paths map[string]string
}

func (fixture *Fixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(fixture.Paths[r.URL.Path])
	if err != nil {
		fixture.t.Fatal(err)
	}

	if _, err := io.Copy(w, f); err != nil {
		fixture.t.Fatal(err)
	}
}

func TestClient_ResourceShow(t *testing.T) {
	mock := httptest.NewServer(
		&Fixture{t,
			map[string]string{
				"/api/3/action/resource_show": "./test/resource.json",
			},
		},
	)
	BaseURL = mock.URL

	client := NewClient()
	actual, err := client.ResourceShow(context.Background(), "blah-blah")
	require.NoError(t, err)

	timestamp, err := time.Parse(ResourceCreatedTimeFormat, "2019-11-12 01:00:00")
	require.NoError(t, err)

	expected := Resource{
		ID:           "1235678-1234-1234-1234-000123456789",
		LastModified: LastModifiedTime{timestamp.Add(12 * time.Hour)},
		Revisions: []Revision{
			{
				ID:              "12112019_2",
				MimeType:        "application/json",
				Name:            "example.json",
				Format:          "JSON",
				URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_2",
				ResourceCreated: ResourceCreatedTime{timestamp.Add(12 * time.Hour)},
				Size:            10000000,
			},
			{
				ID:              "12112019_1",
				MimeType:        "application/json",
				Name:            "example.json",
				Format:          "JSON",
				URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_1",
				ResourceCreated: ResourceCreatedTime{timestamp},
				Size:            20000000,
			},
		},
		PackageID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, expected, *actual)
}

func TestClient_ResourceRevision(t *testing.T) {
	mock := httptest.NewServer(
		&Fixture{t,
			map[string]string{
				"/dataset/00000000-0000-0000-0000-000000000000/resource/235678-1234-1234-1234-000123456789/revision/12112019_1": "./test/example.json",
			},
		},
	)
	BaseURL = mock.URL

	client := NewClient()
	_, err := client.ResourceRevision(context.Background(),
		"00000000-0000-0000-0000-000000000000",
		"235678-1234-1234-1234-000123456789",
		"12112019_1",
	)
	assert.NoError(t, err)
}
