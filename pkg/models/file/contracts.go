package file

import (
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
)

type ContractModel struct {
	DB *[]models.Contract
}

// df := dataframe.LoadStructs(contracts)
// fmt.Println(df.Col("Rohmarge").Quantile(0.1))
//groups := df.GroupBy("In_area")
//fmt.Println(groups)

func (c *ContractModel) Index() (*[]models.Contract, error) {
	return c.DB, nil
}

func (c *ContractModel) Get(id int) (*models.Contract, error) {
	return &(*c.DB)[id], nil
}
