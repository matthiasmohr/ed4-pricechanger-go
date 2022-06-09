package main

import (
	"errors"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/data"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"net/http"
	"time"
)

func (app *application) indexContractsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductSerialNumber string
		ProductNames        []string
		NewPriceInclude     string
		data.Filters
	}

	//v := validator.New()
	qs := r.URL.Query()

	// Read and define the query parameters
	input.ProductSerialNumber = app.readString(qs, "ProductSerialNumber", "")
	input.ProductNames = app.readCSV(qs, "ProductNames", []string{})
	input.NewPriceInclude = app.readString(qs, "NewPriceInclude", "")
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "ProductSerialNumber")
	input.Filters.SortSafelist = []string{"ProductNames", "-ProductNames", "ProductSerialNumber", "-ProductSerialNumber"}

	/*
	   if data.ValidateFilters(v, input.Filters); !v.Valid() {
	       app.failedValidationResponse(w, r, v.Errors)
	       return
	   }
	*/

	c, metadata, err := app.contracts.Index(input.ProductSerialNumber, input.ProductNames, input.NewPriceInclude, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"contracts": c, "metadata": metadata}, nil)
	if err != nil {
		app.errorLog.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showContractHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDStringParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	c, err := app.contracts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"contract": c}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) editContractHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDStringParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	c, err := app.contracts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
	}

	var input struct {
		NewPriceInclude   bool    `json:"NewPriceInclude"`
		NewPriceBase      float64 `json:"NewPriceBase"`
		NewPriceKwh       float64 `json:"NewPriceKwh"`
		NewPriceStartdate string  `json:"NewPriceStartdate"`
	}

	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the request body to the appropriate fields of the contract
	c.NewPriceInclude = input.NewPriceInclude
	c.NewPriceBase = input.NewPriceBase
	c.NewPriceKwh = input.NewPriceKwh
	c.NewPriceStartdate = input.NewPriceStartdate

	cNew, err := app.contracts.Update(c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"contract": cNew}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) editContractsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductSerialNumber string
		ProductNames        []string
		data.Filters
		Typeofchange string
		Change       string
		Changebase   float64
		Changekwh    float64
		Changedate   string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Adjust Input
	input.Filters.Page = 1
	input.Filters.PageSize = 9999999

	contracts, metadata, err := app.contracts.Index(input.ProductSerialNumber, input.ProductNames, "", input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	adjusted := 0
	for _, c := range *contracts {
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
		case "date":
			switch input.Change {
			case "take":
				// TODO: Berechnen
				c.NewPriceStartdate = time.Now().Local().AddDate(0, 2, 0).Format("2006-01-02")
			case "set":
				c.NewPriceStartdate = input.Changedate
			case "add":
				// TODO
			default:
				app.serverErrorResponse(w, r, nil)
				return
			}
		case "communication":
			switch input.Change {
			default:
				app.serverErrorResponse(w, r, nil)
				return
				// TODO
			}
		default:
			app.serverErrorResponse(w, r, nil)
			return
		}

		_, err := app.contracts.Update(&c)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		adjusted = adjusted + 1
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "adjusted": adjusted}, nil)
	if err != nil {
		app.errorLog.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

// --------- STATISTICS HANDLERS -------------

func (app *application) aggregateContractsHandler(w http.ResponseWriter, r *http.Request) {
	aggregator, err := app.readIDStringParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	c, t, err := app.contracts.Aggregate(aggregator)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{aggregator: c, aggregator + "Transposed": t}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) describeContractsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.Describe()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"describeContracts": c}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) quantileContractsHandler(w http.ResponseWriter, r *http.Request) {
	kpi, err := app.readIDStringParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	c, err := app.contracts.Quantile(100, kpi)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"quantile": c}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
