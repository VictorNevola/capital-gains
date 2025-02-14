package strategies

import (
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/pkg/infrastructure/serialization"
)

type BuyStrategy struct{}

func NewBuyStrategy() *BuyStrategy {
	return &BuyStrategy{}
}

func (s *BuyStrategy) CalculateTax(operation entities.Operation, taxState *entities.TaxState) entities.TaxResult {
	taxState.CalculateWeightedMean(operation)
	taxState.AddStockQuantities(operation.StockQuantity)

	return entities.TaxResult{
		Tax: serialization.KeepZero(entities.MinTaxFee),
	}
}
