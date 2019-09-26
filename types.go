package govdata

import "encoding/json"

// Response is a general structure of platform response.
type Response struct {
	Result json.RawMessage `json:"result"`
}

// Resource represents detailed information about resource and it's changes.
type Resource struct {
	Revisions []Revision `json:"resource_revisions"`
	PackageID string     `json:"package_id"`
}

// Revision is an represents changes of a resource.
type Revision struct {
	ID              string  `json:"id"`
	MimeType        string  `json:"mimetype"`
	Name            string  `json:"name"`
	Format          string  `json:"format"`
	URL             string  `json:"url"`
	FileHashSum     *string `json:"file_hash_sum"`
	ResourceCreated string  `json:"resource_created"`
	Size            int     `json:"size"`
}
