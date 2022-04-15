package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

type Contract struct {
	Product_name string  `csv:"product_name"`
	In_area      string  `csv:"in_area"`
	Rohmarge     float64 `csv:"Rohmarge/KWH"`
}
