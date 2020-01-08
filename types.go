package govdata

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Response is a general structure of platform response.
type Response struct {
	Result json.RawMessage `json:"result"`
}

// Package represents detailed information about package and it's changes.
type Package struct {
	Resources []Resource `json:"resources"`
}

// Resource represents detailed information about resource and it's changes.
type Resource struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Revisions []Revision `json:"resource_revisions"`
	PackageID string     `json:"package_id"`
}

// TimeFormat is default time format of the government website.
const TimeFormat = "2006-01-02 15:04:05"

// Time is almost the same as time.Time, but has behaves differently on JSON serialization.
type Time struct {
	time.Time
}

// UnmarshalJSON overrides JSON deserialization.
func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(TimeFormat, s)
	return
}

// MarshalJSON overrides JSON serialization.
func (ct *Time) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == (time.Time{}).UnixNano() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(TimeFormat))), nil
}

// Revision is an represents changes of a resource.
type Revision struct {
	ID              string  `json:"id"`
	MimeType        string  `json:"mimetype"`
	Name            string  `json:"name"`
	Format          string  `json:"format"`
	URL             string  `json:"url"`
	FileHashSum     *string `json:"file_hash_sum"`
	ResourceCreated Time    `json:"resource_created"`
	Size            int     `json:"size"`
}
