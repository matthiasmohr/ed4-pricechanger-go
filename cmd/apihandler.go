package main

import (
	"errors"
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
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
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

func (app *application) showChartHandler(w http.ResponseWriter, r *http.Request) {
	c, t, err := app.contracts.AnalyseByProducts()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"AnalyseByProducts": c, "AnalyseByProductsTransposed": t}, nil)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
