package model

type Schemas struct {
	Instance Schema `json:"instance"`
	// There can be non-bindable services, hence 'omitempty'
	Binding Schema `json:"binding,omitempty"`
}

// A schema consists of the schema for inputs and the schema for outputs.
// Schemas are in the form of JSON Schema v4 (http://json-schema.org/).
type Schema struct {
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}
