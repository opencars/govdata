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
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Revisions    []Revision       `json:"resource_revisions"`
	PackageID    string           `json:"package_id"`
	URL          string           `json:"url"`
	LastModified LastModifiedTime `json:"last_modified"`
}

// Revision is an represents changes of a resource.
type Revision struct {
	ResourceID      string              `json:"resource_id"`
	ID              string              `json:"id"`
	MimeType        string              `json:"mimetype"`
	Name            string              `json:"name"`
	Format          string              `json:"format"`
	URL             string              `json:"url"`
	FileHashSum     *string             `json:"file_hash_sum"`
	ResourceCreated ResourceCreatedTime `json:"resource_created"`
	Size            int                 `json:"size"`
}

// TimeFormat is default time format of the government website.
const (
	ResourceCreatedTimeFormat = "2006-01-02 15:04:05"
	LastModifiedTimeFormat    = "2006-01-02T15:04:05.999999"
)

type (
	// ResourceCreatedTime is almost the same as time.Time, but has behaves differently on JSON serialization.
	ResourceCreatedTime struct{ time.Time }
	// LastModifiedTime is almost the same as time.Time, but has behaves differently on JSON serialization.
	LastModifiedTime struct{ time.Time }
)

// UnmarshalJSON overrides JSON deserialization.
func (ct *ResourceCreatedTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ResourceCreatedTimeFormat, s)
	return
}

// MarshalJSON overrides JSON serialization.
func (ct *ResourceCreatedTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == (time.Time{}).UnixNano() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ResourceCreatedTimeFormat))), nil
}

// UnmarshalJSON overrides JSON deserialization.
func (ct *LastModifiedTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(LastModifiedTimeFormat, s)
	return
}

// MarshalJSON overrides JSON serialization.
func (ct *LastModifiedTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == (time.Time{}).UnixNano() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(LastModifiedTimeFormat))), nil
}
