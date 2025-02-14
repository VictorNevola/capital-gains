package tests_test

import (
	"capital-gains-cli/internal/domain/entities"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateWeightedMean(t *testing.T) {
	taxState := &entities.TaxState{
		WeightedMean:    10.0,
		StockQuantities: 100.0,
	}

	operation := entities.Operation{
		UnitCost:      20.0,
		StockQuantity: 50.0,
	}

	taxState.CalculateWeightedMean(operation)

	expectedWeightedMean := ((100.0 * 10.0) + (50.0 * 20.0)) / (100.0 + 50.0)
	expectedWeightedMean = math.Round(expectedWeightedMean*100) / 100

	assert.Equal(t, expectedWeightedMean, taxState.WeightedMean)
}

func TestAddStockQuantities(t *testing.T) {
	taxState := &entities.TaxState{
		StockQuantities: 100.0,
	}

	taxState.AddStockQuantities(50.0)

	assert.Equal(t, 150.0, taxState.StockQuantities)
}

func TestRemoveStockQuantities(t *testing.T) {
	taxState := &entities.TaxState{
		StockQuantities: 100.0,
	}

	taxState.RemoveStockQuantites(50.0)

	assert.Equal(t, 50.0, taxState.StockQuantities)

	taxState.RemoveStockQuantites(60.0)

	assert.Equal(t, 0.0, taxState.StockQuantities)
}

func TestSetCumulativeLoss(t *testing.T) {
	taxState := &entities.TaxState{
		CumulativeLoss: 100.0,
	}

	taxState.SetCumulativeLoss(-50.0)

	assert.Equal(t, 150.0, taxState.CumulativeLoss)

	taxState.SetCumulativeLoss(30.0)

	assert.Equal(t, 120.0, taxState.CumulativeLoss)

	taxState.SetCumulativeLoss(150.0)

	assert.Equal(t, 0.0, taxState.CumulativeLoss)
}
