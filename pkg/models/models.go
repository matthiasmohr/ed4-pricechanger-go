package models

type Contract struct {
	Product_name string  `csv:"product_name"`
	In_area      string  `csv:"in_area"`
	Rohmarge     float64 `csv:"Rohmarge/KWH"`
}