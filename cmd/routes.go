package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	// JSON response
	//mux.HandleFunc("/v1/products", app.createKreditangebotHandler).Methods("POST")
	mux.HandleFunc("/v1/contracts", app.indexContractsHandler)
	//mux.HandleFunc("/v1/tools/new", app.createToolForm)
	//mux.HandleFunc("/v1/tools/{id}", app.deleteTool).Methods("DELETE")
	mux.HandleFunc("/v1/contract/{id}", app.showContractHandler).Methods("GET")
	mux.HandleFunc("/v1/contract/{id}", app.editContractHandler).Methods("PUT")
	//mux.HandleFunc("/v1/tools/list", app.indexJSON)
	mux.HandleFunc("/v1/aggregateContracts", app.aggregateHandler).Methods("GET")

	// WEB Response
	mux.HandleFunc("/", app.homeWebHandler)
	mux.HandleFunc("/produktverteilung", app.produktverteilungWebHandler)
	mux.HandleFunc("/renditenhistogramm", app.renditenhistogrammWebHandler)
	mux.HandleFunc("/renditevspreis", app.renditevspreisWebHandler)
	mux.HandleFunc("/renditevslaufzeit", app.renditevslaufzeitWebHandler)

	// Misc
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public", fileServer))

	return mux
}

// Convert the notFoundResponse() helper to a http.Handler using the
// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
// Not Found responses.
// router.NotFound = http.HandlerFunc(app.notFoundResponse)

// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
// it as the custom error handler for 405 Method Not Allowed responses.
// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
