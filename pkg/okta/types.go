package okta

// Schema represents Okta user profile scheme.
type Schema struct {
	ID          string                    `json:"id,omitempty"`
	Schema      string                    `json:"$schema,omitempty"`
	Name        string                    `json:"name,omitempty"`
	Title       string                    `json:"title,omitempty"`
	LastUpdated string                    `json:"lastUpdated,omitempty"`
	Created     string                    `json:"created,omitempty"`
	Definitions map[string]Definition     `json:"definitions,omitempty"`
	Type        string                    `json:"type,omitempty"`
	Properties  map[string]SchemaProperty `json:"properties,omitempty"`
}

// Definition represents Okta definition.
type Definition struct {
	ID         string                        `json:"id,omitempty"`
	Type       string                        `json:"type,omitempty"`
	Properties map[string]DefinitionProperty `json:"properties,omitempty"`
	Required   []string                      `json:"required,omitempty"`
}

// DefinitionProperty represents Okta definition property.
type DefinitionProperty struct {
	Title     string `json:"title,omitempty"`
	Type      string `json:"type,omitempty"`
	Required  *bool  `json:"required,omitempty"`
	Scope     string `json:"scope,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`
}

// SchemaProperty represents Okta scheme property.
type SchemaProperty struct {
	AllOf []Ref `json:"allOf,omitempty"`
}

// Ref represents Okta scheme reference.
type Ref struct {
	Ref string `json:"ref,omitempty"`
}

// User represents user structure.
type User struct {
	ID     string
	Login  string
	Status string
}

// Group represents user group structure.
type Group struct {
	ID          string
	Name        string
	Description string
}
