package db

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/gocarina/gocsv"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/data"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type ContractModel struct {
	DB *gorm.DB
}

func (c *ContractModel) Index(ProductSerialNumber string, ProductNames []string, NewPriceInclude string, Commodity string, filters data.Filters) (*[]models.Contract, data.Metadata, error) {
	var contracts = []models.Contract{}

	// State base query
	result := c.DB.
		Limit(filters.PageSize).
		Offset(filters.PageSize*filters.Page - filters.PageSize).
		Order("product_serial_number").
		Where(&models.Contract{ProductSerialNumber: ProductSerialNumber}).
		Where(&models.Contract{Commodity: Commodity})
	// Iterate ProductNames in Query
	if len(ProductNames) > 0 {
		result = result.Where("product_name IN ?", ProductNames)
	}
	// Check if included in new prices
	switch NewPriceInclude {
	case "false":
		result = result.Where(map[string]interface{}{"new_price_include": false})
		break
	case "true":
		result = result.Where(&models.Contract{NewPriceInclude: true})
		break
	}

	// Execute query and check for errors
	result = result.Find(&contracts)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, data.Metadata{}, result.Error
	}

	// Count the numbers
	var amount int64
	count := c.DB.Model(&models.Contract{}).
		Where(&models.Contract{ProductSerialNumber: ProductSerialNumber}).
		Where(&models.Contract{Commodity: Commodity})
	// Iterate ProductNames in Query
	if len(ProductNames) > 0 {
		count = count.Where("product_name IN ?", ProductNames)
	}
	switch NewPriceInclude {
	case "false":
		count = count.Where(map[string]interface{}{"new_price_include": false})
		break
	case "true":
		count = count.Where(&models.Contract{NewPriceInclude: true})
		break
	}
	result = result.Count(&amount)

	metadata := data.CalculateMetadata(int(amount), filters.Page, filters.PageSize)

	return &contracts, metadata, nil
}

func (c *ContractModel) Get(id string) (*models.Contract, error) {
	var contract = models.Contract{}

	if result := c.DB.First(&contract, "product_serial_number = ?", id); result.Error != nil {
		return nil, models.ErrNoRecord
	}

	return &contract, nil
}

func (c *ContractModel) Update(contract *models.Contract) (*models.Contract, error) {
	result := c.DB.Save(contract)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	return contract, nil
}

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

// ---- DATABASE -----

func Init() *gorm.DB {
	dsn := "host=localhost user=postgres password=q1alfa147 dbname=ed4-pricechanger port=5432 sslmode=disable TimeZone=Europe/Berlin"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.Contract{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func (c *ContractModel) Reset(env string) error {
	c.DB.Migrator().DropTable(&models.Contract{})
	c.DB.AutoMigrate(&models.Contract{})

	var filename string
	// Open Database Connection
	if env == "production" {
		filename = "dataInput/20220531 Stage Data.csv"
	} else {
		filename = "dataInput/20220525 frontEndPayload.csv"
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Parse File into Contract Slice and Calculate Total Prices
	var contractDB []models.Contract
	if err := gocsv.UnmarshalFile(f, &contractDB); err != nil {
		return err
	}
	for i, _ := range contractDB {
		contractDB[i].CalculateTotalPrices()
	}

	for i, _ := range contractDB {
		if result := c.DB.Create(&contractDB[i]); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
