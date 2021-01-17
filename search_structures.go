package main

// Filter search
type Filter struct {
	Name    []FieldValue `json:"name,omitempty"`
	Artist  []FieldValue `json:"artist,omitempty"`
	Album   []FieldValue `json:"album,omitempty"`
	Genre   []FieldValue `json:"genre,omitempty"`
	Year    []FieldValue `json:"year,omitempty"`
	Country []FieldValue `json:"country,omitempty"`
	Offset  string       `json:"offset,omitempty"`
	Type    string       `json:"type,omitempty"`
}

// Paging response object ...
type Paging struct {
	Href     string      `json:"href"`
	Items    interface{} `json:"items"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous string      `json:"previous"`
	Total    int         `json:"total"`
}

// FieldValue doc ...
type FieldValue struct {
	Value   string
	Exclude bool
}

// FieldFilter  ...
type FieldFilter struct {
	Name   string
	Values []FieldValue
}
