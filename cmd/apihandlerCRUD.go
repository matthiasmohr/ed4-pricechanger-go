package main

import (
	"errors"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/data"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"net/http"
)

func (app *application) indexContractsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductSerialNumber string
		ProductNames        []string
		NewPriceInclude     string
		Commodity           string
		data.Filters
	}

	//v := validator.New()
	qs := r.URL.Query()

	// Read and define the query parameters
	input.ProductSerialNumber = app.readString(qs, "ProductSerialNumber", "")
	input.ProductNames = app.readCSV(qs, "ProductNames", []string{})
	input.NewPriceInclude = app.readString(qs, "NewPriceInclude", "")
	input.Commodity = app.readString(qs, "Commodity", "")
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

	c, metadata, err := app.contracts.Index(input.ProductSerialNumber, input.ProductNames, input.NewPriceInclude, input.Commodity, input.Filters)
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

// --------- DATABASE  -------------
func (app *application) databaseReset(w http.ResponseWriter, r *http.Request) {
	err := app.contracts.Reset(app.config.env)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Reset executed. Errors:": err}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
