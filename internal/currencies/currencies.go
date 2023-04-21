package Currency

type MinMaxCurrency struct {
	Value float64
	Date  string
}
type AvgCurrency struct {
	Value    float64
	Sum      float64
	Quantity int
}

type CurrRate struct {
	Date       string `xml:"Date,attr"`
	Currencies []struct {
		Name  string `xml:"Name"`
		Value string `xml:"Value"`
	} `xml:"Valute"`
}
