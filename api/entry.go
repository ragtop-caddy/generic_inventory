package main

// Transaction - Standard transaction structure
type Transaction struct {
	ISODate  string `json:"isodate,omitempty"`
	Campus   string `json:"campus,omitempty"`
	Students int    `json:"students,omitempty"`
	Action   string `json:"action,omitempty"`
	Amount   int    `json:"amount,omitempty"`
}

// Header - Information standard to all entries
type Header struct {
	Type        string        `json:"type,omitempty"`
	Description string        `json:"description,omitempty"`
	Stock       int           `json:"stock,omitempty"`
	History     []Transaction `json:"history,omitempty"`
}

// Detail - Defines entry details
type Detail struct {
	Gender string `json:"gender,omitempty"`
	Color  string `json:"color,omitempty"`
	Size   string `json:"size,omitempty"`
	Style  string `json:"style,omitempty"`
	Fit    string `json:"fit,omitempty"`
}

// Entry - Defines various types of inventory
type Entry struct {
	SKU     string  `json:"sku,omitempty"`
	Header  *Header `json:"header,omitempty"`
	Details *Detail `json:"details,omitempty"`
}

var entries []Entry
