package main

import (
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"net/http"
	"time"
)

func (app *application) editContractsHandler(w http.ResponseWriter, r *http.Request) {
	input := models.Pricechangerequest

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Adjust Input
	input.Filters.Page = 1
	input.Filters.PageSize = 99999999

	// Get List of Products
	contracts, metadata, err := app.contracts.Index(input.ProductSerialNumber, input.ProductNames, "", input.Commodity, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//var wg sync.WaitGroup
	adjusted := 0
	for _, c := range *contracts {
		// Check for Exclusion reasons
		if c.StartDate < input.ExcludeProductchangeFrom {
			c.NewPriceInclude = false
			continue
		}
		// Add termination date comparison
		// Add Product Change date comparison
		switch input.Typeofchange {
		case "price":
			switch input.Change {
			case "take":
				c.NewPriceInclude = true
				c.NewPriceBase = c.BaseNewPriceProposed
				c.NewPriceKwh = c.KwhNewPriceProposed
			case "set":
				c.NewPriceInclude = true
				c.NewPriceBase = input.Changebase
				c.NewPriceKwh = input.Changekwh
			case "add":
				c.NewPriceInclude = true
				c.NewPriceBase = c.NewPriceBase + input.Changebase
				c.NewPriceKwh = c.NewPriceKwh + input.Changekwh
			case "exclude":
				c.NewPriceInclude = false
			default:
				app.serverErrorResponse(w, r, nil)
				return
			}

			// Check for the Limits
			if input.LimitToCatalogueprice == true {
				if c.CatalogKwhPrice < c.NewPriceKwh {
					c.NewPriceKwh = c.CatalogKwhPrice
				}
				if c.CatalogBasePrice < c.NewPriceBase {
					c.NewPriceBase = c.CatalogBasePrice
				}
			}
			if input.LimitToMax == true {
				if c.NewPriceKwh > input.LimitToMaxKwh {
					c.NewPriceKwh = input.LimitToMaxKwh
				}
				if c.NewPriceBase > input.LimitToMaxBase {
					c.NewPriceBase = input.LimitToMaxBase
				}
			}
			if input.LimitToMin == true {
				if c.NewPriceKwh < input.LimitToMinKwh {
					c.NewPriceKwh = input.LimitToMinKwh
				}
				if c.NewPriceBase < input.LimitToMinBase {
					c.NewPriceBase = input.LimitToMinBase
				}
			}
			if input.LimitToFactor == true {
				if c.NewPriceKwh > (c.CurrentKwhPriceNet * float64(input.LimitToFactorMax/100)) {
					c.NewPriceKwh = c.CurrentKwhPriceNet * float64(input.LimitToFactorMax/100)
				}
				if c.NewPriceBase > (c.CurrentBasePriceNet * float64(input.LimitToFactorMax/100)) {
					c.NewPriceBase = c.CurrentBasePriceNet * float64(input.LimitToFactorMax/100)
				}
				if c.NewPriceKwh < (c.CurrentKwhPriceNet * float64(input.LimitToFactorMin/100)) {
					c.NewPriceKwh = c.CurrentKwhPriceNet * float64(input.LimitToFactorMin/100)
				}
				if c.NewPriceBase < (c.CurrentBasePriceNet * float64(input.LimitToFactorMin/100)) {
					c.NewPriceBase = c.CurrentBasePriceNet * float64(input.LimitToFactorMin/100)
				}
			}
			if input.LimitToContractduration == true {
				if input.LimitToContractdurationDays > 1 {
					c.NewPriceInclude = false
					// TODO: Calculate contract duration
				}
			}
		case "date":
			switch input.Change {
			case "take":
				// TODO: Berechnen
				c.NewPriceStartdate = time.Now().Local().AddDate(0, 2, 0).Format("2006-01-02")
			case "set":
				c.NewPriceStartdate = input.Changedate
			default:
				app.serverErrorResponse(w, r, nil)
				return
			}
		case "communication":
			switch input.Channel {
			case "postmail":
				c.CommunicationChannel = "postmail"
				c.CalculateCommunicationDates(input.Allatonce, input.AllatonceDate, input.Beforechange, input.BeforechangeDays)
			case "email":
				c.CommunicationChannel = "email"
				c.CalculateCommunicationDates(input.Allatonce, input.AllatonceDate, input.Beforechange, input.BeforechangeDays)
			case "both":
				c.CommunicationChannel = "both"
				c.CalculateCommunicationDates(input.Allatonce, input.AllatonceDate, input.Beforechange, input.BeforechangeDays)
			case "none":
				c.CommunicationChannel = ""
				c.CommunicationDate1 = ""
				c.CommunicationDate2 = ""
			default:
				app.serverErrorResponse(w, r, nil)
				return
			}
		default:
			app.serverErrorResponse(w, r, nil)
			return
		}

		//wg.Add(1)
		d := c // Create a copy of c in order to not hit a race condition (it is only a pointer)
		go func() {
			//defer wg.Done()
			_, error := app.contracts.Update(&d)
			if error != nil {
				app.logError(nil, err)
			}
		}()
		adjusted = adjusted + 1
	}
	//wg.Wait()

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "adjusted": adjusted}, nil)
	if err != nil {
		app.errorLog.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
