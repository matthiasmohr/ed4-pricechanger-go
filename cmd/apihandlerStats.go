package main

import (
	"errors"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"net/http"
)

// --------- STATISTICS HANDLERS -------------

func (app *application) aggregateContractsHandler(w http.ResponseWriter, r *http.Request) {
	aggregator, err := app.readIDStringParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Read params
	qs := r.URL.Query()
	groupby := app.readString(qs, "groupby", "ProductName")
	commodity := app.readString(qs, "Commodity", "")

	c, t, err := app.contracts.Aggregate(groupby, aggregator, commodity)
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
	// Read params
	qs := r.URL.Query()
	commodity := app.readString(qs, "Commodity", "")

	c, err := app.contracts.Describe(commodity)
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

	// Read params
	qs := r.URL.Query()
	commodity := app.readString(qs, "Commodity", "")

	c, err := app.contracts.Quantile(100, kpi, commodity)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"quantile": c}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
