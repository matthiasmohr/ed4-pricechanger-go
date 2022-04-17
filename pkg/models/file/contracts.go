package file

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
)

type ContractModel struct {
	DB *[]models.Contract
}

func (c *ContractModel) Index() (*[]models.Contract, error) {
	return c.DB, nil
}

func (c *ContractModel) Get(id int) (*models.Contract, error) {
	return &(*c.DB)[id], nil
}

func (c *ContractModel) AnalyseByProducts() (map[string][]interface{}, error) {
	df := dataframe.LoadStructs(*c.DB)
	//fmt.Println(df.Col("Rohmarge").Quantile(0.1))
	groups := df.GroupBy("Product_name")
	aggre := groups.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_MAX, dataframe.Aggregation_MEAN, dataframe.Aggregation_MIN, dataframe.Aggregation_COUNT}, []string{"Rohmarge", "Rohmarge", "Rohmarge", "Rohmarge"}) // Maximum value in column "values",  Minimum value in column "values2"
	records := aggre.Records()

	// Transpose records and convert to map of Type
	// map[Product_name:[ELECTRICITY_SUBSCRIPTION_24 ELECTRICITY_SUBSCRIPTION_12] Rohmarge_COUNT:[3.000000 1497.000000] Rohmarge_MAX:[0.147761 26.200000] Rohmarge_MEAN:[0.142167 0.171183] Rohmarge_MIN:[0.136481 0.108396]]
	resultVectors := make(map[string][]interface{})
	for i := 0; i < len(records[0]); i++ {
		for k := 0; k < len(records); k++ {
			if k != 0 {
				resultVectors[records[0][i]] = append(resultVectors[records[0][i]], records[k][i])
			}
		}
	}

	// TODO: Error Catching
	return resultVectors, nil
}
