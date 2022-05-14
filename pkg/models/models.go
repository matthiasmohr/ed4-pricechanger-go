package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

type Contract struct {
	ProductName         string // `csv:"product_name"`
	InArea              string
	MbaId               string
	ProductSerialNumber string
	commodity           string

	// Contract terms
	AnnualConsumption float64
	StartDate         string //date
	ContractDuration  string //date
	OrderDate         string //date
	Status            string
	PriceGuarantee    string //date

	// Original pricing data
	PriceDate           string //date
	PriceChange         bool
	OrigBaseCosts       float64
	OrigKwhCosts        float64
	OrigKwhMargin       float64
	OrigbaseMargin      float64
	CurrentBasePriceNet float64 // as calculated from origCost & origMargin
	CurrentKwhPriceNet  float64 // as calculated from origCost & origMargin

	// New price data
	CatalogBasePrice     float64 //price that would be applicable for new signup
	CatalogKwhPrice      float64
	CurrentBaseCost      float64 //from current cost data (ENET or B7)
	CurrentKwhCost       float64
	CurrentKwhRohmarge   float64
	CurrentBaseRohmarge  float64
	BaseNewPriceProposed float64 //currentCost + origMargin
	KwhNewPriceProposed  float64 //currentCost + origMargin

	// Price Change info
	NewPriceInclude   bool
	NewPriceBase      float64
	NewPriceKwh       float64
	NewPriceStartdate string // date
}
