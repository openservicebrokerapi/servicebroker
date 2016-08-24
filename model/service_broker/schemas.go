package model

type Schemas struct {
	Instance Schema `json:"instance"`
	Binding  Schema `json:"binding"`
}

// A schema consists of the schema for inputs and the schema for outputs.
// Schemas are in the form of JSON Schema v4 (http://json-schema.org/).
type Schema struct {
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}
