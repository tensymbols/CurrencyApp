package main

import (
	"ValuteApp/internal/app"
	"ValuteApp/internal/date"
	"ValuteApp/internal/ports"
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

	valuteApp := app.NewValuteApp()
	client := ports.NewClient()

	var wg sync.WaitGroup
	doneCh := make(chan any)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < days; i++ {

			currDate := dateutil.GetStringDate(date)

			dateutil.DateDecrease(&date)
			resp, err := client.GetWithParam(url, "date_req", currDate)
			if err != nil {
				client.AddError(err)
			}
			byteValCurs, err := client.GetData(resp)
			if err != nil {
				client.AddError(err)
			}

			currValCurs, err := valuteApp.ParseValCurs(byteValCurs)

			if err != nil {
				valuteApp.AddError(err)
			}

			valuteApp.AddValCurs(currValCurs)

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

	valuteApp.ProcessAll()

	valuteApp.PrintMinCurs()
	valuteApp.PrintMaxCurs()
	valuteApp.PrintAvgCurs()

	fmt.Println(client.CheckErrors())
	fmt.Println(valuteApp.CheckErrors())

}
