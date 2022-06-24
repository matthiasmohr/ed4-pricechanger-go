package models

import (
	"gorm.io/gorm"
)

type Adjustment struct {
	gorm.Model `dataframe:"-"`

	ProductName string `gorm:"index"`
	InArea      bool
	Commodity   string `gorm:"index"`

	// Contract terms
	AdjustmentKwh  float64
	AdjustmentBase float64
}
