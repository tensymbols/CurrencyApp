package main

import (
	"CurrencyCB/internal/app"
	"CurrencyCB/internal/date"
	"CurrencyCB/internal/ports"
	"flag"
	"fmt"
	"sync"
	"time"
)

func main() {

	var days int
	url := "http://www.cbr.ru/scripts/XML_daily_eng.asp"
	date := time.Now()
	flag.IntVar(&days, "days", 90, "number of days")
	flag.Parse()

	CurrencyApp := app.NewCurrencyApp()
	client := ports.NewClient()

	var wg sync.WaitGroup
	doneCh := make(chan any)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < days; i++ {

			tempDate := dateutil.GetStringDate(date)

			dateutil.DateDecrease(&date)
			resp, err := client.GetWithParam(url, "date_req", tempDate)
			if err != nil {
				client.AddError(err)
			}
			byteCurrRate, err := client.GetData(resp)
			if err != nil {
				client.AddError(err)
			}

			tempCurrRate, err := CurrencyApp.ParseCurrRate(byteCurrRate)

			if err != nil {
				CurrencyApp.AddError(err)
			}

			CurrencyApp.AddCurrRate(tempCurrRate)

			resp.Body.Close()
		}
		doneCh <- true
	}()
	go func() {
		fmt.Print("\rЗагрузка...")
		for range doneCh {
			fmt.Printf("\r\n")
			return
		}

	}()
	wg.Wait()

	CurrencyApp.ProcessAll()

	CurrencyApp.PrintMinRate()
	CurrencyApp.PrintMaxRate()
	CurrencyApp.PrintAvgRate()

	fmt.Println(client.CheckErrors())
	fmt.Println(CurrencyApp.CheckErrors())

}
