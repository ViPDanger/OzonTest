package dto

import "encoding/xml"

// Valcursdto represents the root element of the XML answer with currency courses.
// @description The structure of the response with foreign exchange rates for the specified date and name.
type ValCursDTO struct {
	XMLName xml.Name    `xml:"ValCurs" swaggerignore:"true"`
	Date    string      `xml:"Date,attr" example:"02.01.2006"`              // Source Date
	Name    string      `xml:"name,attr" example:"Foreign Currency Market"` // Source name
	Valutes []ValuteDTO `xml:"Valute"`                                      // List of currencies
}

// Valutedto represents a separate currency in the list of courses.
// @description information about the currency and its course.
type ValuteDTO struct {
	ID        string `xml:"ID,attr" example:"R01235"`    // Currency identifier
	NumCode   string `xml:"NumCode" example:"840"`       // Numerical currency code (ISO 4217)
	CharCode  string `xml:"CharCode" example:"USD"`      // Lamp currency code (ISO 4217)
	Nominal   string `xml:"Nominal" example:"1"`         // Nominal
	Name      string `xml:"Name" example:"Доллар США"`   // Name of the currency
	Value     string `xml:"Value" example:"30,9436"`     // Course value
	VunitRate string `xml:"VunitRate" example:"30,9436"` // Course per unit of currency
}
