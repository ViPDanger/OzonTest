package dto

import "encoding/xml"

type ValCursDTO struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Date    string      `xml:"Date,attr"`
	Name    string      `xml:"name,attr"`
	Valutes []ValuteDTO `xml:"Valute"`
}

type ValuteDTO struct {
	ID        string `xml:"ID,attr"`
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   string `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`     // можно преобразовать в float64 с заменой запятой
	VunitRate string `xml:"VunitRate"` // аналогично
}
