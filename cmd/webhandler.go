package main

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("ui/html/*"))

func (app *application) homeWebHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.Index()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "home.page.tmpl", c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) produktverteilungWebHandler(w http.ResponseWriter, r *http.Request) {
	c, err := app.contracts.AnalyseByProducts()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "produktverteilung.page.tmpl", c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) renditenhistogrammWebHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "renditenhistogramm.page.tmpl", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) renditevspreisWebHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "renditevspreis.page.tmpl", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) renditevslaufzeitWebHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "renditevslaufzeit.page.tmpl", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
