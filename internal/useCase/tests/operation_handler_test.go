package tests_test

import (
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/domain/interfaces/mocks"
	usecase "capital-gains-cli/internal/useCase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOperationHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBuyStrategy := mocks.NewMockTaxCalculator(ctrl)
	mockedSellStrategy := mocks.NewMockTaxCalculator(ctrl)

	operationHandler := usecase.NewOperationHandler(usecase.MapStrategies{
		entities.OperationTypeBuy:  mockedBuyStrategy,
		entities.OperationTypeSell: mockedSellStrategy,
	})

	t.Run("should call the correct strategies with expected parameters", func(t *testing.T) {
		mockedBuyStrategy.EXPECT().CalculateTax(
			gomock.Eq(entities.Operation{
				Operation:     entities.OperationTypeBuy,
				UnitCost:      10.00,
				StockQuantity: 10000,
			}),
			gomock.Any(),
		).Return(entities.TaxResult{})

		mockedSellStrategy.EXPECT().CalculateTax(
			gomock.Eq(entities.Operation{
				Operation:     entities.OperationTypeSell,
				UnitCost:      5.00,
				StockQuantity: 5000,
			}),
			gomock.Any(),
		).Return(entities.TaxResult{})

		operationHandler.Handler([][]entities.Operation{
			{
				{
					Operation:     entities.OperationTypeBuy,
					UnitCost:      10.00,
					StockQuantity: 10000,
				},
			},
			{
				{
					Operation:     entities.OperationTypeSell,
					UnitCost:      5.00,
					StockQuantity: 5000,
				},
			},
		})
	})

	t.Run("should return the correct tax results", func(t *testing.T) {
		mockedBuyStrategy.EXPECT().CalculateTax(gomock.Any(), gomock.Any()).Return(entities.TaxResult{
			Tax: 0.00,
		})
		mockedSellStrategy.EXPECT().CalculateTax(gomock.Any(), gomock.Any()).Return(entities.TaxResult{
			Tax: 100.00,
		})

		results, err := operationHandler.Handler([][]entities.Operation{
			{
				{
					Operation:     entities.OperationTypeBuy,
					UnitCost:      10.00,
					StockQuantity: 10000,
				},
			},
			{
				{
					Operation:     entities.OperationTypeSell,
					UnitCost:      5.00,
					StockQuantity: 5000,
				},
			},
		})

		assert.NoError(t, err)
		assert.EqualValues(t, [][]entities.TaxResult{
			{{Tax: 0.00}},
			{{Tax: 100.00}},
		}, results)
	})

	t.Run("should return an error if the operation type is not supported", func(t *testing.T) {
		_, err := operationHandler.Handler([][]entities.Operation{
			{{Operation: "invalid"}},
		})
		assert.Error(t, err)
		assert.ErrorIs(t, err, usecase.ErrOperationTypeNotSupported)
	})
}
