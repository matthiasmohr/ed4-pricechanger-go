package models

import "github.com/matthiasmohr/ed4-pricechanger-go/internal/data"

var Pricechangerequest struct {
	data.Filters
	Typeofchange                string
	Change                      string
	Changebase                  float64
	Changekwh                   float64
	ProductNames                []string
	ProductSerialNumber         string
	Commodity                   string
	ExcludeSOS                  bool
	ExcludeSOSFrom              string
	ExcludeContractduration     bool
	ExcludeContractdurationDays int
	ExcludeTermination          bool
	ExcludeTerminationFrom      string
	ExcludeProductchange        bool
	ExcludeProductchangeFrom    string
	LimitToCatalogueprice       bool
	LimitToMax                  bool
	LimitToMaxBase              float64
	LimitToMaxKwh               float64
	LimitToMin                  bool
	LimitToMinBase              float64
	LimitToMinKwh               float64
	LimitToFactor               bool
	LimitToFactorMin            float64
	LimitToFactorMax            float64

	Changedate string

	Channel          string
	Allatonce        bool
	AllatonceDate    string
	Beforechange     bool
	BeforechangeDays int
}
