package govdata

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Fixture struct {
	t    *testing.T
	Path string
}

func (fixture *Fixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(fixture.Path)
	if err != nil {
		fixture.t.Fatal(err)
	}

	if _, err := io.Copy(w, f); err != nil {
		fixture.t.Fatal(err)
	}
}

func TestClient_ResourceShow(t *testing.T) {
	mock := httptest.NewServer(&Fixture{t, "./test/resource.json"})
	BaseURL = mock.URL

	client := NewClient()
	actual, err := client.ResourceShow(context.Background(), "blah-blah")
	assert.NoError(t, err)

	expected := Resource{
		Revisions: []Revision{
			{
				ID:              "12112019_2",
				MimeType:        "application/json",
				Name:            "example.json",
				Format:          "JSON",
				URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_2",
				ResourceCreated: "2019-11-12 12:00:00",
				Size:            10000000,
			},
			{
				ID:              "12112019_1",
				MimeType:        "application/json",
				Name:            "example.json",
				Format:          "JSON",
				URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_1",
				ResourceCreated: "2019-11-12 01:00:00",
				Size:            20000000,
			},
		},
		PackageID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, expected, *actual)
}

func TestClient_ResourceRevision(t *testing.T) {
	mock := httptest.NewServer(&Fixture{t, "./test/example.json"})
	BaseURL = mock.URL

	client := NewClient()
	_, err := client.ResourceRevision(context.Background(), "", "", "")
	assert.NoError(t, err)
}
