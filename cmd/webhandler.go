package main

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))

func (app *application) indexWebHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.Index()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "Index", c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) graphWebHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.AnalyseByProducts()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "Graph", c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
