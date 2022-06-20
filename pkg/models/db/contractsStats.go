package db

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/data"
)

// ---- STATISTICS -----

func (c *ContractModel) Aggregate(groupby string, aggregator string, commodity string) ([]map[string]interface{}, map[string][]interface{}, error) {
	allcontracts, _, _ := c.Index("", nil, "", commodity, data.Filters{PageSize: 999999})
	df := dataframe.LoadStructs(*allcontracts)
	groups := df.GroupBy(groupby)
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

func (c *ContractModel) Describe(commodity string) ([]map[string]interface{}, error) {
	allcontracts, _, _ := c.Index("", nil, "", commodity, data.Filters{PageSize: 999999})
	df := dataframe.LoadStructs(*allcontracts)
	dfDescribe := df.Describe()
	outputMap := dfDescribe.Maps()

	// TODO: Error Catching
	return outputMap, nil
}

func (c *ContractModel) Quantile(n int, column string, commodity string) ([]float64, error) {
	allcontracts, _, _ := c.Index("", nil, "", commodity, data.Filters{PageSize: 999999})
	df := dataframe.LoadStructs(*allcontracts)
	var array []float64
	for i := 1; i < n; i++ {
		quantile := df.Col(column).Quantile(float64(i) / float64(n))
		array = append(array, quantile)
	}

	// TODO: Error Catching
	return array, nil
}
