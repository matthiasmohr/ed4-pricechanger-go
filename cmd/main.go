package main

import (
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models/file"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	contracts *file.ContractModel
	config    config
}

type config struct {
	port int
	env  string
}

func main() {
	var cfg config

	// Load Flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Define Logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open Database Connection
	f, err := os.Open("data_afsr.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var contractDB []models.Contract
	if err := gocsv.UnmarshalFile(f, &contractDB); err != nil {
		panic(err)
	}

	// Start App
	app := &application{
		config:    cfg,
		errorLog:  errorLog,
		infoLog:   infoLog,
		contracts: &file.ContractModel{DB: &contractDB},
	}

	// Start Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
