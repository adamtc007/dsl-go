package manager

// Attribute represents a single entry in the data dictionary.
type Attribute struct {
	AttributeID string `json:"AttributeID"`
	Description string `json:"Description"`
	VectorID    string `json:"VectorID"`
}
