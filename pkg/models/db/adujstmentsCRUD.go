package db

import (
	"fmt"
	"github.com/matthiasmohr/ed4-pricechanger-go/pkg/models"
	"gorm.io/gorm"
)

type AdjustmentModel struct {
	DB *gorm.DB
}

func (a *AdjustmentModel) Index() (*[]models.Adjustment, error) {
	var adjustments = []models.Adjustment{}

	// State base query
	result := a.DB.Find(&adjustments)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return &adjustments, nil
}

func (a *AdjustmentModel) Get(id string) (*models.Adjustment, error) {
	var adjustment = models.Adjustment{}

	if result := a.DB.First(&adjustment, "id = ?", id); result.Error != nil {
		return nil, models.ErrNoRecord
	}

	return &adjustment, nil
}

func (a *AdjustmentModel) Update(adjustment *models.Adjustment) (*models.Adjustment, error) {
	result := a.DB.Save(adjustment)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	return adjustment, nil
}

// ---- DATABASE -----

func (a *AdjustmentModel) Reset(env string) error {
	a.DB.Migrator().DropTable(&models.Adjustment{})
	a.DB.AutoMigrate(&models.Adjustment{})
	return nil
}
