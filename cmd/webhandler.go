package main

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("ui/html/*"))

func (app *application) indexWebHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.Index()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "list.page.tmpl", c)
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

	err = tmpl.ExecuteTemplate(w, "graph.page.tmpl", c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) graph2WebHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "graph2.page.tmpl", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
