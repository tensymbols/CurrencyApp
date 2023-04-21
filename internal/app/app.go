package app

import (
	"ValuteApp/internal/valutes"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const line = "------------------------------------------------------------------"

type ValuteError struct {
	Err error
}

type ValuteErrors []ValuteError

func (ve ValuteError) Error() string {
	return ve.Err.Error()
}

func (ve ValuteErrors) Error() string {
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

type ValuteApp struct {
	valCurses    []valutes.ValCurs //количество дней
	valuteErrors ValuteErrors

	minCurs map[string]valutes.MinMaxValute
	maxCurs map[string]valutes.MinMaxValute
	avgCurs map[string]valutes.AvgValute
}

func NewValuteApp() ValuteApp {
	return ValuteApp{
		minCurs: map[string]valutes.MinMaxValute{},
		maxCurs: map[string]valutes.MinMaxValute{},
		avgCurs: map[string]valutes.AvgValute{},
	}
}

func (va *ValuteApp) CheckErrors() string {
	if len(va.valuteErrors) > 0 {
		return va.valuteErrors.Error()
	} else {
		return "No Valute App errors"
	}
}

func (va *ValuteApp) AddError(err error) {
	va.valuteErrors = append(va.valuteErrors, err.(ValuteError))
}

func (va *ValuteApp) parseCursVal(val string) (float64, error) {
	strVal := strings.Replace(val, ",", ".", -1)
	fVal, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		return 0, ValuteError{fmt.Errorf("error parsing float value from %s"+err.Error(), val)}
	}
	return fVal, nil
}
func (va *ValuteApp) AddValCurs(curs valutes.ValCurs) {
	va.valCurses = append(va.valCurses, curs)
}

func (va *ValuteApp) ParseValCurs(buf []byte) (valutes.ValCurs, error) {
	vc := valutes.ValCurs{}

	dec := xml.NewDecoder(bytes.NewReader(buf))
	dec.CharsetReader = func(enc string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	err := dec.Decode(&vc)
	if err != nil {
		return valutes.ValCurs{}, ValuteError{err}
	}
	return vc, nil
}

func (va *ValuteApp) ProcessAll() {
	err := va.processMinCurs()
	if err != nil {
		va.AddError(err)
	}
	err = va.processMaxCurs()
	if err != nil {
		va.AddError(err)
	}
	err = va.processAvgCurs()
	if err != nil {
		va.AddError(err)
	}
}

func (va *ValuteApp) processMinCurs() error {
	for _, curs := range va.valCurses {
		for _, val := range curs.Valute {
			fVal, err := va.parseCursVal(val.Value)
			if err != nil {
				return err
			}
			v, ok := va.minCurs[val.Name]
			if !ok {
				va.minCurs[val.Name] = valutes.MinMaxValute{Value: fVal, Date: curs.Date}
			} else if fVal < v.Value {
				va.minCurs[val.Name] = valutes.MinMaxValute{Value: fVal, Date: curs.Date}
			}
		}
	}

	return nil
}

func (va *ValuteApp) processMaxCurs() error {
	for _, curs := range va.valCurses {
		for _, val := range curs.Valute {
			fVal, err := va.parseCursVal(val.Value)
			if err != nil {
				return err
			}
			v, ok := va.maxCurs[val.Name]
			if !ok {
				va.maxCurs[val.Name] = valutes.MinMaxValute{Value: fVal, Date: curs.Date}
			} else if fVal < v.Value {
				va.maxCurs[val.Name] = valutes.MinMaxValute{Value: fVal, Date: curs.Date}
			}
		}
	}
	return nil
}

func (va *ValuteApp) processAvgCurs() error {
	for _, curs := range va.valCurses {
		for _, val := range curs.Valute {
			fVal, err := va.parseCursVal(val.Value)
			if err != nil {
				return err
			}
			v, ok := va.avgCurs[val.Name]
			if !ok {
				va.avgCurs[val.Name] = valutes.AvgValute{Sum: fVal, Quantity: 1}
			} else if fVal < v.Value {
				va.avgCurs[val.Name] = valutes.AvgValute{Sum: v.Sum + fVal, Quantity: v.Quantity + 1}
			}
		}
	}
	for k, v := range va.avgCurs {
		va.avgCurs[k] = valutes.AvgValute{Sum: v.Sum, Quantity: v.Quantity, Value: v.Sum / float64(v.Quantity)}
	}
	return nil
}

func (va *ValuteApp) PrintMinCurs() {
	fmt.Println("\rМинимальные курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ | ДАТА")

	for k, v := range va.minCurs {
		fmt.Printf(" %s | %.2f | %s\n", k, v.Value, v.Date)
	}

	fmt.Println(line)
}
func (va *ValuteApp) PrintMaxCurs() {
	fmt.Println("\rМаксимальные курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ | ДАТА")

	for k, v := range va.maxCurs {
		fmt.Printf(" %s | %.2f | %s\n", k, v.Value, v.Date)
	}
	fmt.Println(line)
}
func (va *ValuteApp) PrintAvgCurs() {
	fmt.Println("\rСредние курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ")

	for k, v := range va.avgCurs {
		fmt.Printf(" %s | %.2f\n", k, v.Value)
	}
	fmt.Println(line)
}
