package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Contract struct {
	gorm.Model `dataframe:"-"`

	ProductName         string `gorm:"index"`
	InArea              bool
	MbaId               string
	ProductSerialNumber string `gorm:"unique, index"`
	Commodity           string `gorm:"index" csv:"commodity"`

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

	// Total Prices
	OrigTotalCosts        float64
	CurrentTotalPriceNet  float64
	TotalNewPriceProposed float64

	// Price Change info
	NewPriceInclude      bool `gorm:"index"`
	NewPriceBase         float64
	NewPriceKwh          float64
	NewPriceTotal        float64
	NewPriceStartdate    string `gorm:"index"` // date
	CommunicationChannel string
	CommunicationDate1   string
	CommunicationDate2   string
}

func (c *Contract) CalculateTotalPrices() {
	c.OrigTotalCosts = c.OrigBaseCosts + c.OrigKwhCosts/100*c.AnnualConsumption
	c.CurrentTotalPriceNet = c.CurrentBasePriceNet + c.CurrentKwhPriceNet/100*c.AnnualConsumption
	c.TotalNewPriceProposed = c.BaseNewPriceProposed + c.KwhNewPriceProposed/100*c.AnnualConsumption
	c.NewPriceTotal = c.NewPriceBase + c.NewPriceKwh/100*c.AnnualConsumption
}

func (c *Contract) CalculateCommunicationDates(Allatonce bool, AllatonceDate string, Beforechange bool, Beforechangedays int) {
	NewPriceStartdateDate, err := time.Parse("2006-01-02", c.NewPriceStartdate)
	if err != nil {
		fmt.Println("Fehler bei der Datumsberechnung: ", err)
	}
	Beforechangedate := NewPriceStartdateDate.AddDate(0, 0, -Beforechangedays)
	if Allatonce && Beforechange {
		c.CommunicationDate1 = AllatonceDate
		c.CommunicationDate2 = Beforechangedate.Format("2006-01-02")
	} else if Allatonce {
		c.CommunicationDate1 = AllatonceDate
	} else if Beforechange {
		c.CommunicationDate1 = Beforechangedate.Format("2006-01-02")
	}
}
