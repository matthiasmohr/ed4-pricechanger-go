package file

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
)

type ContractModel struct {
	DB *[]models.Contract
}

func (c *ContractModel) Index() (*[]models.Contract, error) {
	var result []models.Contract
	for i, v := range *c.DB {
		if i < 100 {
			result = append(result, v)
		}
	}
	return &result, nil
}

func (c *ContractModel) Get(id string) (*models.Contract, error) {
	for _, v := range *c.DB {
		if v.ProductSerialNumber == id {
			return &v, nil
		}
	}
	return nil, models.ErrNoRecord
}

func (c *ContractModel) Update(d *models.Contract) (*models.Contract, error) {
	for _, v := range *c.DB {
		if v.ProductSerialNumber == d.ProductSerialNumber {
			v.ProductName = "AAAAAAAAAAA"
			return &v, nil
		}
	}
	return nil, models.ErrNoRecord
}

func (c *ContractModel) Aggregate(aggregator string) ([]map[string]interface{}, map[string][]interface{}, error) {
	df := dataframe.LoadStructs(*c.DB)
	//fmt.Println(df.Col("Rohmarge").Quantile(0.1))
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
	// map[Product_name:[ELECTRICITY_SUBSCRIPTION_24 ELECTRICITY_SUBSCRIPTION_12] Rohmarge_COUNT:[3.000000 1497.000000] Rohmarge_MAX:[0.147761 26.200000] Rohmarge_MEAN:[0.142167 0.171183] Rohmarge_MIN:[0.136481 0.108396]]
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
	df := dataframe.LoadStructs(*c.DB)
	dfDescribe := df.Describe()
	outputMap := dfDescribe.Maps()

	// TODO: Error Catching
	return outputMap, nil
}

func (c *ContractModel) Quantile(n int, column string) ([]float64, error) {
	df := dataframe.LoadStructs(*c.DB)
	var array []float64
	for i := 1; i < n; i++ {
		quantile := df.Col(column).Quantile(1 / float64(n))
		array = append(array, quantile)
	}

	// TODO: Error Catching
	return array, nil
}
