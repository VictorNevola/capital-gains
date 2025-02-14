package tests_test

import (
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/domain/strategies"
	"capital-gains-cli/internal/pkg/infrastructure/serialization"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSellStrategy(t *testing.T) {
	t.Parallel()

	sellStrategy := strategies.SellStrategy{
		FeeTax:      0.20,
		AmountToTax: 20000,
	}

	t.Run("should return tax fee 0 when selling stocks and amount operation is less than amount to tax", func(t *testing.T) {
		t.Parallel()

		taxState := &entities.TaxState{}
		operation := entities.Operation{
			StockQuantity: 50,
			UnitCost:      15.00,
		}

		taxFee := sellStrategy.CalculateTax(operation, taxState)

		assert.Equal(t, serialization.KeepZero(0.00), taxFee.Tax)
		assert.Equal(t, 0.00, taxState.StockQuantities)
	})

	t.Run("should return tax fee 0 when selling stocks and have loss in operation", func(t *testing.T) {
		t.Parallel()

		taxState := &entities.TaxState{
			StockQuantities: 5000,
			WeightedMean:    10.00,
			CumulativeLoss:  0,
		}

		operation := entities.Operation{
			StockQuantity: 5000,
			UnitCost:      5.00,
		}

		taxFee := sellStrategy.CalculateTax(operation, taxState)

		assert.Equal(t, serialization.KeepZero(0.00), taxFee.Tax)
		assert.Equal(t, 0.00, taxState.StockQuantities)
		assert.Equal(t, 10.00, taxState.WeightedMean)
		assert.Equal(t, 25000.00, taxState.CumulativeLoss)
	})

	t.Run("should calculate tax fee when profit is above tax threshold and there is no previous loss", func(t *testing.T) {
		t.Parallel()

		taxState := &entities.TaxState{
			StockQuantities: 10000,
			WeightedMean:    10.00,
			CumulativeLoss:  0,
		}

		operation := entities.Operation{
			StockQuantity: 5000,
			UnitCost:      20.00,
		}

		taxFee := sellStrategy.CalculateTax(operation, taxState)

		assert.Equal(t, serialization.KeepZero(10000.00), taxFee.Tax)
		assert.Equal(t, 5000.00, taxState.StockQuantities)
		assert.Equal(t, 10.00, taxState.WeightedMean)
		assert.Equal(t, 0.00, taxState.CumulativeLoss)
	})

	t.Run("should calculate tax fee when profit is above tax threshold and there is previous loss", func(t *testing.T) {
		t.Parallel()

		taxState := &entities.TaxState{
			StockQuantities: 5000,
			WeightedMean:    10.00,
			CumulativeLoss:  25000,
		}

		operation := entities.Operation{
			StockQuantity: 3000,
			UnitCost:      20.00,
		}

		taxFee := sellStrategy.CalculateTax(operation, taxState)

		assert.Equal(t, serialization.KeepZero(1000.00), taxFee.Tax)
		assert.Equal(t, 2000.00, taxState.StockQuantities)
		assert.Equal(t, 10.00, taxState.WeightedMean)
		assert.Equal(t, 0.00, taxState.CumulativeLoss)
	})
}
