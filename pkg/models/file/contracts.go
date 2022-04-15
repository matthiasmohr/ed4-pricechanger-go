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

func (m *ContractModel) Index() ([]*models.Contracts, error){
	stmt := "SELECT * FROM kreditangebote ORDER BY id DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kreditangebote := []*models.Kreditangebot{}

	for rows.Next() {
		k := &models.Kreditangebot{}
		err := rows.Scan(&k.Id, &k.SollZins, &k.EffektivZins, &k.Zinsbindung, &k.Abfragedatum)
		if err != nil {
			return nil, err
		}
		kreditangebote = append(kreditangebote, k)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return kreditangebote, nil
}

func (m *KreditangebotModel) Get(id int) (*models.Kreditangebot, error){
	stmt := "SELECT * FROM kreditangebote WHERE id=?"
	row := m.DB.QueryRow(stmt, id)
	k := &models.Kreditangebot{}
	err := row.Scan(&k.Id, &k.SollZins, &k.EffektivZins, &k.Zinsbindung, &k.Abfragedatum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return k, nil
}