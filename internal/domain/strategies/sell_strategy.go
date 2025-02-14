package strategies

import (
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/pkg/infrastructure/serialization"
	"capital-gains-cli/internal/pkg/utils"
)

type SellStrategy struct {
	FeeTax      float64
	AmountToTax float64
}

func NewSellStrategy(feeTax float64, amountToTax float64) *SellStrategy {
	return &SellStrategy{
		FeeTax:      feeTax,
		AmountToTax: amountToTax,
	}
}

func (s *SellStrategy) CalculateTax(operation entities.Operation, taxState *entities.TaxState) entities.TaxResult {
	operationProfit := s.calculatProfit(operation, taxState)
	taxFee := s.calculateFeeTax(
		operation,
		operationProfit,
		taxState.CumulativeLoss,
	)

	taxState.SetCumulativeLoss(operationProfit)
	taxState.RemoveStockQuantites(operation.StockQuantity)

	return entities.TaxResult{
		Tax: serialization.KeepZero(taxFee),
	}
}

func (s *SellStrategy) calculatProfit(operation entities.Operation, taxState *entities.TaxState) float64 {
	return (operation.UnitCost - taxState.WeightedMean) * operation.StockQuantity
}

func (s *SellStrategy) calculateFeeTax(
	operation entities.Operation,
	operationProfit float64,
	cumulativeLoss float64,
) float64 {
	operationAmount := operation.UnitCost * operation.StockQuantity
	operationProfit = operationProfit - cumulativeLoss

	if operationAmount <= s.AmountToTax || operationProfit <= entities.MinOperationProfit {
		return entities.MinTaxFee
	}

	taxFee := operationProfit * s.FeeTax

	return utils.RoundToTwoDecimalPlaces(taxFee)
}
