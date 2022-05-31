package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()
	mux.Use(CORS)

	// JSON response
	//mux.HandleFunc("/v1/products", app.createKreditangebotHandler).Methods("POST")
	mux.HandleFunc("/v1/contracts", app.indexContractsHandler)
	mux.HandleFunc("/v1/updatecontracts", app.editContractsHandler).Methods("PUT", "OPTIONS")
	//mux.HandleFunc("/v1/tools/new", app.createToolForm)
	//mux.HandleFunc("/v1/tools/{id}", app.deleteTool).Methods("DELETE")
	mux.HandleFunc("/v1/contract/{id}", app.showContractHandler).Methods("GET")
	mux.HandleFunc("/v1/contract/{id}", app.editContractHandler).Methods("PUT", "OPTIONS")
	//mux.HandleFunc("/v1/tools/list", app.indexJSON)
	mux.HandleFunc("/v1/describeContracts", app.describeHandler).Methods("GET")
	mux.HandleFunc("/v1/aggregateContracts/{id}", app.aggregateHandler).Methods("GET")
	mux.HandleFunc("/v1/quantileContracts/{id}", app.quantileHandler).Methods("GET")

	// Misc
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public", fileServer))

	return mux
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

// Convert the notFoundResponse() helper to a http.Handler using the
// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
// Not Found responses.
// router.NotFound = http.HandlerFunc(app.notFoundResponse)

// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
// it as the custom error handler for 405 Method Not Allowed responses.
// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
