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

func (va *ValuteApp) parseCursVal(val string) (float32, error) {
	strVal := strings.Replace(val, ",", ".", -1)
	fVal64, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		return 0, ValuteError{fmt.Errorf("error parsing float value from %s"+err.Error(), val)}
	}
	fVal := float32(fVal64)
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
	for _, curs := range va.valCurses {
		for _, val := range curs.Valute {

			err := va.processMinCurs(val.Name, val.Value, curs.Date)
			if err != nil {
				va.AddError(err)
			}
			err = va.processMaxCurs(val.Name, val.Value, curs.Date)
			if err != nil {
				va.AddError(err)
			}
			err = va.processAvgCurs(val.Name, val.Value)
			if err != nil {
				va.AddError(err)
			}

		}
	}
}

func (va *ValuteApp) processMinCurs(name string, val string, date string) error {

	fVal, err := va.parseCursVal(val)
	if err != nil {
		return err
	}
	v, ok := va.minCurs[name]
	if !ok {
		va.minCurs[name] = valutes.MinMaxValute{Value: fVal, Date: date}
	} else if fVal < v.Value {
		va.minCurs[name] = valutes.MinMaxValute{Value: fVal, Date: date}
	}
	return nil
}

func (va *ValuteApp) processMaxCurs(name string, val string, date string) error {

	fVal, err := va.parseCursVal(val)
	if err != nil {
		return err
	}
	v, ok := va.maxCurs[name]
	if !ok {
		va.maxCurs[name] = valutes.MinMaxValute{Value: fVal, Date: date}
	} else if fVal > v.Value {
		va.maxCurs[name] = valutes.MinMaxValute{Value: fVal, Date: date}
	}
	return nil
}

func (va *ValuteApp) processAvgCurs(name string, val string) error {
	fVal, err := va.parseCursVal(val)
	if err != nil {
		return ValuteError{err}
	}
	v, ok := va.avgCurs[name]
	if !ok {
		va.avgCurs[name] = valutes.AvgValute{Value: fVal, Quantity: 1}
	} else {
		newAvg := ((float32(v.Quantity) * v.Value) + fVal) / (float32(v.Quantity + 1))
		va.avgCurs[name] = valutes.AvgValute{Value: newAvg, Quantity: v.Quantity + 1}
	}
	return nil
}

func (va *ValuteApp) PrintMinCurs() {
	fmt.Println("\rМинимальные курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ | ДАТА")

	for k, v := range va.minCurs {
		fmt.Println(k, "|", v.Value, "|", v.Date)
	}

	fmt.Println(line)
}
func (va *ValuteApp) PrintMaxCurs() {
	fmt.Println("\rМаксимальные курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ | ДАТА")

	for k, v := range va.maxCurs {
		fmt.Println(k, "|", v.Value, "|", v.Date)
	}
	fmt.Println(line)
}
func (va *ValuteApp) PrintAvgCurs() {
	fmt.Println("\rСредние курсы валют:\n" +
		"ИМЯ ВАЛЮТЫ | ЗНАЧЕНИЕ")

	for k, v := range va.avgCurs {
		fmt.Println(k, "|", v.Value)
	}
	fmt.Println(line)
}
