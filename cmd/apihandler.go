package main

import (
	"errors"
	"fmt"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"net/http"
)

func (app *application) indexContractsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.Index()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"contracts": c}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
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
	fmt.Println(err)
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

func (app *application) aggregateHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) describeHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) quantileHandler(w http.ResponseWriter, r *http.Request) {
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
