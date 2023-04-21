package app

import (
	"CurrencyCB/internal/currencies"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const line = "------------------------------------------------------------------"

type CurrencyError struct {
	Err error
}

type CurrencyErrors []CurrencyError

func (ve CurrencyError) Error() string {
	return ve.Err.Error()
}

func (ve CurrencyErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var es string
	es += ve[0].Error()
	for i := 1; i < len(ve); i++ {
		es += "\n" + ve[i].Error()
	}
	return es
}

type CurrencyApp struct {
	CurrencyErrors CurrencyErrors

	currRates []Currency.CurrRate //количество дней

	minRate map[string]Currency.MinMaxCurrency
	maxRate map[string]Currency.MinMaxCurrency
	avgRate map[string]Currency.AvgCurrency
}

func NewCurrencyApp() CurrencyApp {
	return CurrencyApp{
		minRate: map[string]Currency.MinMaxCurrency{},
		maxRate: map[string]Currency.MinMaxCurrency{},
		avgRate: map[string]Currency.AvgCurrency{},
	}
}

func (va *CurrencyApp) CheckErrors() string {
	if len(va.CurrencyErrors) > 0 {
		return va.CurrencyErrors.Error()
	} else {
		return "No Currencies App errors"
	}
}

func (va *CurrencyApp) AddError(err error) {
	va.CurrencyErrors = append(va.CurrencyErrors, err.(CurrencyError))
}

func (va *CurrencyApp) parseRateVal(val string) (float64, error) {
	strVal := strings.Replace(val, ",", ".", -1)
	fVal, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		return 0, CurrencyError{fmt.Errorf("error parsing float value from %s"+err.Error(), val)}
	}
	return fVal, nil
}
func (va *CurrencyApp) AddCurrRate(rate Currency.CurrRate) {
	va.currRates = append(va.currRates, rate)
}

func (va *CurrencyApp) ParseCurrRate(buf []byte) (Currency.CurrRate, error) {
	vc := Currency.CurrRate{}

	dec := xml.NewDecoder(bytes.NewReader(buf))
	dec.CharsetReader = func(enc string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	err := dec.Decode(&vc)
	if err != nil {
		return Currency.CurrRate{}, CurrencyError{err}
	}
	return vc, nil
}

func (va *CurrencyApp) ProcessAll() {
	err := va.processMinRate()
	if err != nil {
		va.AddError(err)
	}
	err = va.processMaxRate()
	if err != nil {
		va.AddError(err)
	}
	err = va.processAvgRate()
	if err != nil {
		va.AddError(err)
	}
}

func (va *CurrencyApp) processMinRate() error {
	for _, rate := range va.currRates {
		for _, val := range rate.Currencies {
			fVal, err := va.parseRateVal(val.Value)
			if err != nil {
				return err
			}
			v, ok := va.minRate[val.Name]
			if !ok {
				va.minRate[val.Name] = Currency.MinMaxCurrency{Value: fVal, Date: rate.Date}
			} else if fVal < v.Value {
				va.minRate[val.Name] = Currency.MinMaxCurrency{Value: fVal, Date: rate.Date}
			}
		}
	}

	return nil
}

func (va *CurrencyApp) processMaxRate() error {
	for _, rate := range va.currRates {
		for _, val := range rate.Currencies {
			fVal, err := va.parseRateVal(val.Value)
			if err != nil {
				return err
			}
			v, ok := va.maxRate[val.Name]
			if !ok {
				va.maxRate[val.Name] = Currency.MinMaxCurrency{Value: fVal, Date: rate.Date}
			} else if fVal < v.Value {
				va.maxRate[val.Name] = Currency.MinMaxCurrency{Value: fVal, Date: rate.Date}
			}
		}
	}
	return nil
}

func (va *CurrencyApp) processAvgRate() error {
	for _, rate := range va.currRates {
		for _, val := range rate.Currencies {
			fVal, err := va.parseRateVal(val.Value)
			if err != nil {
				return err
			}
			v, ok := va.avgRate[val.Name]
			if !ok {
				va.avgRate[val.Name] = Currency.AvgCurrency{Sum: fVal, Quantity: 1}
			} else if fVal < v.Value {
				va.avgRate[val.Name] = Currency.AvgCurrency{Sum: v.Sum + fVal, Quantity: v.Quantity + 1}
			}
		}
	}
	for k, v := range va.avgRate {
		va.avgRate[k] = Currency.AvgCurrency{Sum: v.Sum, Quantity: v.Quantity, Value: v.Sum / float64(v.Quantity)}
	}
	return nil
}

func (va *CurrencyApp) PrintMinRate() {
	fmt.Println("\rМинимальные курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ | ДАТА")

	for k, v := range va.minRate {
		fmt.Printf(" %s | %.2f | %s\n", k, v.Value, v.Date)
	}

	fmt.Println(line)
}
func (va *CurrencyApp) PrintMaxRate() {
	fmt.Println("\rМаксимальные курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ | ДАТА")

	for k, v := range va.maxRate {
		fmt.Printf(" %s | %.2f | %s\n", k, v.Value, v.Date)
	}
	fmt.Println(line)
}
func (va *CurrencyApp) PrintAvgRate() {
	fmt.Println("\rСредние курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ")

	for k, v := range va.avgRate {
		fmt.Printf(" %s | %.2f\n", k, v.Value)
	}
	fmt.Println(line)
}
