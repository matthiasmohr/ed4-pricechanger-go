package db

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/data"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

type ContractModel struct {
	mutex sync.Mutex
	DB    *gorm.DB
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
	if len(ProductNames) > 0 && len(ProductNames[0]) > 0 {
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
	// Mutex für parallele Bearbeitung
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result := c.DB.Save(contract)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	return contract, nil
}

// ---- DATABASE -----

func Init(env string) *gorm.DB {
	var dsn string
	if env == "production" {
		dsn = "host=ec2-52-72-56-59.compute-1.amazonaws.com user=bkctbapxeeddmq password=73309a32852b2457f074bb14e65a28d67239a762471da1053bba59a0c912cfa3 dbname=ddp7q4qn5eseu9 port=5432 sslmode=require TimeZone=Europe/Berlin"
	} else {
		dsn = "host=localhost user=postgres password=q1alfa147 dbname=ed4-pricechanger port=5432 sslmode=disable TimeZone=Europe/Berlin"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	sqlDB, err := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(1)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(2)

	//Test the database
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.Contract{}, &models.Adjustment{})
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
