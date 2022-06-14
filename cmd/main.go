package main

import (
	"flag"
	"fmt"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models/db"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	contracts *db.ContractModel
	db        *gorm.DB
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

	var database = db.Init(cfg.env)

	// Start App
	app := &application{
		config:   cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
		//contracts: &file.ContractModel{DB: contractDB},
		contracts: &db.ContractModel{DB: database},
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

	infoLog.Printf("Starting server on %s in %s", srv.Addr, cfg.env)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
