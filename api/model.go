package api

// Transaction - Standard transaction structure
type Transaction struct {
	ISODate  string `bson:"isodate,omitempty" json:"isodate,omitempty"`
	Campus   string `bson:"campus,omitempty" json:"campus,omitempty"`
	Students int    `bson:"students,omitempty" json:"students,omitempty"`
	Action   string `bson:"action,omitempty" json:"action,omitempty"`
	Amount   int    `bson:"amount,omitempty" json:"amount,omitempty"`
}

// Header - Information standard to all entries
type Header struct {
	Type        string        `bson:"type,omitempty" json:"type,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Stock       int           `bson:"stock,omitempty" json:"stock,omitempty"`
	History     []Transaction `bson:"history,omitempty" json:"history,omitempty"`
}

// Detail - Defines entry details
type Detail struct {
	Gender string `bson:"gender,omitempty" json:"gender,omitempty"`
	Color  string `bson:"color,omitempty" json:"color,omitempty"`
	Size   string `bson:"size,omitempty" json:"size,omitempty"`
	Style  string `bson:"style,omitempty" json:"style,omitempty"`
	Fit    string `bson:"fit,omitempty" json:"fit,omitempty"`
}

// Entry - Defines various types of inventory
type Entry struct {
	SKU     string  `bson:"sku,omitempty" json:"sku,omitempty"`
	Header  *Header `bson:"header,omitempty" json:"header,omitempty"`
	Details *Detail `bson:"details,omitempty" json:"details,omitempty"`
}
