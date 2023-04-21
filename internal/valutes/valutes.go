package valutes

type MinMaxValute struct {
	Value float64
	Date  string
}
type AvgValute struct {
	Value    float64
	Sum      float64
	Quantity int
}

type ValCurs struct {
	Date   string `xml:"Date,attr"`
	Valute []struct {
		Name  string `xml:"Name"`
		Value string `xml:"Value"`
	} `xml:"Valute"`
}
