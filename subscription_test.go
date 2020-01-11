package govdata

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	mock := httptest.NewServer(
		&Fixture{t,
			map[string]string{
				"/api/3/action/resource_show": "./test/resource.json",
			},
		})
	BaseURL = mock.URL

	timestamp, err := time.Parse(ResourceCreatedTimeFormat, "2019-11-12 01:00:00")
	require.NoError(t, err)

	expected := []Revision{
		{
			ID:              "12112019_2",
			ResourceID:      "1235678-1234-1234-1234-000123456789",
			MimeType:        "application/json",
			Name:            "example.json",
			Format:          "JSON",
			URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_2",
			ResourceCreated: ResourceCreatedTime{timestamp.Add(12 * time.Hour)},
			Size:            10000000,
		},
		{
			ID:              "12112019_1",
			ResourceID:      "1235678-1234-1234-1234-000123456789",
			MimeType:        "application/json",
			Name:            "example.json",
			Format:          "JSON",
			URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_1",
			ResourceCreated: ResourceCreatedTime{timestamp},
			Size:            20000000,
		},
	}

	revisions := Subscribe("1235678-1234-1234-1234-000123456789", timestamp.Add(time.Second))

	select {
	case r := <-revisions:
		assert.Equal(t, expected[0], r)
	case <-time.After(time.Second):
		t.Error("expected revision")
	}
}

func TestSubscribePackage(t *testing.T) {
	mock := httptest.NewServer(
		&Fixture{t,
			map[string]string{
				"/api/3/action/package_show":  "./test/package.json",
				"/api/3/action/resource_show": "./test/resource.json",
			},
		})
	BaseURL = mock.URL

	timestamp, err := time.Parse(ResourceCreatedTimeFormat, "2019-11-12 01:00:00")
	require.NoError(t, err)

	expected := Resource{
		ID:   "1235678-1234-1234-1234-000123456789",
		Name: "example.json",
		URL:  "https://data.gov.ua/dataset/00000000-0000-0000-0000-00000000000/resource/1235678-1234-1234-1234-000123456789/download/example.json",
		Revisions: []Revision{
			{
				ID:              "12112019_2",
				ResourceID:      "1235678-1234-1234-1234-000123456789",
				MimeType:        "application/json",
				Name:            "example.json",
				Format:          "JSON",
				URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_2",
				ResourceCreated: ResourceCreatedTime{timestamp.Add(12 * time.Hour)},
				Size:            10000000,
			},
			{
				ID:              "12112019_1",
				ResourceID:      "1235678-1234-1234-1234-000123456789",
				MimeType:        "application/json",
				Name:            "example.json",
				Format:          "JSON",
				URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_1",
				ResourceCreated: ResourceCreatedTime{timestamp},
				Size:            20000000,
			},
		},
		PackageID:    "00000000-0000-0000-0000-000000000000",
		LastModified: LastModifiedTime{timestamp.Add(12 * time.Hour)},
	}

	events := SubscribePackage("00000000-0000-0000-0000-000000000000", map[string]time.Time{})

	select {
	case r := <-events:
		assert.Equal(t, expected, r)
	case <-time.After(time.Second):
		t.Error("expected revision")
	}
}
