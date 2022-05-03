package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

type Contract struct {
	ProductName         string  `csv:"product_name"`
	InArea              string  `csv:"A"`
	MbaId               string  `csv:"mba_id"`
	ProductSerialNumber string  `csv:"product_serial_number"`
	StartDate           string  `csv:"start_date"`
	LastPriceChange     string  `csv:"last_price_change"`
	BasePriceGross      float64 `csv:"base_price_gross"`
	KwhPriceGross       float64 `csv:"kwh_price_gross"`
	BaseCosts           float64 `csv:"base_costs"`
	KwhCosts            float64 `csv:"kwh_costs"`
	GrossDiscount       float64 `csv:"gross_discount"`
	Consumption         float64 `csv:"consumption"`
	KwhRohmarge         float64 `csv:"kwh_rohmarge"`
	BaseRohmarge        float64 `csv:"base_rohmarge"`
}
