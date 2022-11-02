package data


// Result  Json struct
type Result struct {
	// Meta
	Meta []Metadata `json:"results"`
}

// Metadata metadata
type Metadata struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
	Stars    int    `json:"stars"`
	Version  string `json:"version"`
}


