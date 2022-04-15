package file
import (
	"database/sql"
	"errors"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
)

type ContractModel struct {
	DB *[]models.Contract
}


// df := dataframe.LoadStructs(contracts)
// fmt.Println(df.Col("Rohmarge").Quantile(0.1))
//groups := df.GroupBy("In_area")
//fmt.Println(groups)