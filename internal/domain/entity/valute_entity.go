package entity

type ValuteCurs struct {
	CreatorID string
	Date      string
	Name      string
	Valutes   []Valute
}

type Valute struct {
	ID        string
	NumCode   int
	CharCode  string
	Nominal   int
	Name      string
	Value     float64
	VunitRate float64
}
