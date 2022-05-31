package file

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/data"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
)

type ContractModel struct {
	DB []models.Contract
}

func (c *ContractModel) Index(ProductSerialNumber string, ProductNames []string, filters data.Filters) (*[]models.Contract, data.Metadata, error) {
	var result []models.Contract

	// Filter out the result
	for _, v := range c.DB {
		if ProductSerialNumber == "" || ProductSerialNumber == v.ProductSerialNumber {
			if len(ProductNames) > 0 {
				for _, pn := range ProductNames {
					if v.ProductName == pn || ProductNames == nil {
						result = append(result, v)
					}
				}
			} else {
				result = append(result, v)
			}
		}
	}

	// Sort the result (offset, limit)
	var result2 []models.Contract
	start := filters.PageSize*filters.Page - filters.PageSize
	stop := filters.PageSize * filters.Page
	for i, v := range result {
		if i >= start && i < stop {
			result2 = append(result2, v)
		}
	}

	metadata := data.CalculateMetadata(len(result), filters.Page, filters.PageSize)

	return &result2, metadata, nil
}

func (c *ContractModel) Get(id string) (*models.Contract, error) {
	for _, v := range c.DB {
		if v.ProductSerialNumber == id {
			return &v, nil
		}
	}
	return nil, models.ErrNoRecord
}

func (c *ContractModel) Update(d *models.Contract) (*models.Contract, error) {
	for i, v := range c.DB {
		if v.ProductSerialNumber == d.ProductSerialNumber {
			c.DB[i].NewPriceInclude = d.NewPriceInclude
			c.DB[i].NewPriceBase = d.NewPriceBase
			c.DB[i].NewPriceKwh = d.NewPriceKwh
			c.DB[i].NewPriceStartdate = d.NewPriceStartdate
			c.DB[i].CalculateTotalPrices()
			return &c.DB[i], nil
		}
	}
	return nil, models.ErrNoRecord
}

func (c *ContractModel) Aggregate(aggregator string) ([]map[string]interface{}, map[string][]interface{}, error) {
	df := dataframe.LoadStructs(c.DB)
	groups := df.GroupBy("ProductName")
	aggre := groups.Aggregation([]dataframe.AggregationType{
		dataframe.Aggregation_MAX,
		dataframe.Aggregation_MEAN,
		dataframe.Aggregation_MIN,
		dataframe.Aggregation_COUNT},
		[]string{aggregator, aggregator, aggregator, aggregator})
	outputMap := aggre.Maps()

	// Rename Keys for Client simplicity purpose
	for i, _ := range outputMap {
		outputMap[i]["COUNT"] = outputMap[i][aggregator+"_COUNT"]
		outputMap[i]["MAX"] = outputMap[i][aggregator+"_MAX"]
		outputMap[i]["MEAN"] = outputMap[i][aggregator+"_MEAN"]
		outputMap[i]["MIN"] = outputMap[i][aggregator+"_MIN"]
	}

	// Transpose records and convert to map of Type
	records := aggre.Records()
	transposedVectors := make(map[string][]interface{})
	for i := 0; i < len(records[0]); i++ {
		for k := 0; k < len(records); k++ {
			if k != 0 {
				transposedVectors[records[0][i]] = append(transposedVectors[records[0][i]], records[k][i])
			}
		}
	}

	// TODO: Error Catching
	return outputMap, transposedVectors, nil
}

func (c *ContractModel) Describe() ([]map[string]interface{}, error) {
	df := dataframe.LoadStructs(c.DB)
	dfDescribe := df.Describe()
	outputMap := dfDescribe.Maps()

	// TODO: Error Catching
	return outputMap, nil
}

func (c *ContractModel) Quantile(n int, column string) ([]float64, error) {
	df := dataframe.LoadStructs(c.DB)
	var array []float64
	for i := 1; i < n; i++ {
		quantile := df.Col(column).Quantile(float64(i) / float64(n))
		array = append(array, quantile)
	}

	// TODO: Error Catching
	return array, nil
}
