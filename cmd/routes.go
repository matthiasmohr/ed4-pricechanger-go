package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter().StrictSlash(true)

	// JSON response
	//mux.HandleFunc("/v1/products", app.createKreditangebotHandler).Methods("POST")
	mux.HandleFunc("/v1/contracts", app.indexContractsHandler)
	//mux.HandleFunc("/v1/tools/new", app.createToolForm)
	//mux.HandleFunc("/v1/tools/{id}", app.deleteTool).Methods("DELETE")
	mux.HandleFunc("/v1/contract/{id}", app.showContractHandler).Methods("GET")
	//mux.HandleFunc("/v1/tools/list", app.indexJSON)

	// Misc
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

// Convert the notFoundResponse() helper to a http.Handler using the
// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
// Not Found responses.
// router.NotFound = http.HandlerFunc(app.notFoundResponse)

// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
// it as the custom error handler for 405 Method Not Allowed responses.
// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
