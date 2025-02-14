package tests_test

import (
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/domain/strategies"
	"capital-gains-cli/internal/pkg/infrastructure/serialization"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuyStrategy(t *testing.T) {
	t.Parallel()

	buyStrategy := strategies.BuyStrategy{}

	t.Run("should return tax fee 0 when buying stocks and set new quantity and weighted mean", func(t *testing.T) {
		t.Parallel()
		taxState := &entities.TaxState{}

		operation := entities.Operation{
			StockQuantity: 5,
			UnitCost:      20,
		}

		taxFee := buyStrategy.CalculateTax(operation, taxState)

		assert.Equal(t, serialization.KeepZero(0.00), taxFee.Tax)
		assert.Equal(t, 5.00, taxState.StockQuantities)
		assert.Equal(t, 20.00, taxState.WeightedMean)
	})

	t.Run("should calculate weighted mean as 15.00 when buying 5 stocks at 10.00 with previous 5 stocks at 20.00", func(t *testing.T) {
		t.Parallel()
		taxState := &entities.TaxState{
			StockQuantities: 5,
			WeightedMean:    20,
		}

		operation := entities.Operation{
			StockQuantity: 5,
			UnitCost:      10,
		}

		buyStrategy.CalculateTax(operation, taxState)

		assert.Equal(t, 15.00, taxState.WeightedMean)
	})
}
