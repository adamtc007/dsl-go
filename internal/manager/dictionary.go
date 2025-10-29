package manager

// Attribute represents a single entry in the data dictionary.
type Attribute struct {
	AttributeID string `json:"AttributeID"`
	Description string `json:"Description"`
	VectorID    string `json:"VectorID"`
}

// Product represents a single product in the data dictionary.
type Product struct {
	ProductID   string   `json:"ProductID"`
	Description string   `json:"Description"`
	ServiceIDs  []string `json:"ServiceIDs"`
}

// Service represents a single service in the data dictionary.
type Service struct {
	ServiceID   string `json:"ServiceID"`
	Description string `json:"Description"`
}

// Resource represents a single resource in the data dictionary.
type Resource struct {
	ResourceID  string `json:"ResourceID"`
	Description string `json:"Description"`
}

// DataDictionary represents the entire data dictionary.
type DataDictionary struct {
	Products   []Product   `json:"products"`
	Services   []Service   `json:"services"`
	Resources  []Resource  `json:"resources"`
	Attributes []Attribute `json:"attributes"`
}
