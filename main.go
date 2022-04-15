package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/gocarina/gocsv"
	"log"
	"os"
)

type Contract struct {
	Product_name string  `csv:"product_name"`
	In_area      string  `csv:"in_area"`
	Rohmarge     float64 `csv:"Rohmarge/KWH"`
}

func main() {
	f, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	contracts := []Contract{}

	if err := gocsv.UnmarshalFile(f, &contracts); err != nil {
		panic(err)
	}

	/*for _, contract := range contracts {
		fmt.Println("Hello", contract.Product_name, contract.Rohmarge)
	}*/
	/*
		csvReader := csv.NewReader(f)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
	*/

	df := dataframe.LoadStructs(contracts)
	fmt.Println(df.Col("Rohmarge").Quantile(0.1))
	//groups := df.GroupBy("In_area")
	//fmt.Println(groups)

	// Export CSV File
	// csvContent, err := gocsv.MarshalString(&clients) // Get all clients as CSV string

	// Filter Slice
	// https://stackoverflow.com/questions/67143237/how-to-filter-slice-of-struct-using-filter-parameter-in-go

	// Dataframe Examples
	// https://morioh.com/p/35560c47de92
}
