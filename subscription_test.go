package govdata

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	mock := httptest.NewServer(&Fixture{t, "./test/resource.json"})
	BaseURL = mock.URL

	timestamp, err := time.Parse(TimeFormat, "2019-11-12 01:00:00")
	require.NoError(t, err)

	expected := []Revision{
		{
			ID:              "12112019_2",
			MimeType:        "application/json",
			Name:            "example.json",
			Format:          "JSON",
			URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_2",
			ResourceCreated: Time{timestamp.Add(12 * time.Hour)},
			Size:            10000000,
		},
		{
			ID:              "12112019_1",
			MimeType:        "application/json",
			Name:            "example.json",
			Format:          "JSON",
			URL:             "https://data.gov.ua/dataset/00000000-0000-0000-0000-000000000000/resource/1235678-1234-1234-1234-000123456789/revision/12112019_1",
			ResourceCreated: Time{timestamp},
			Size:            20000000,
		},
	}

	revisions := Subscribe("test", timestamp.Add(time.Second))

	select {
	case r := <-revisions:
		assert.Equal(t, expected[0], r)
	case <-time.After(time.Second):
		t.Error("expected revision")
	}
}
