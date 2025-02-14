package entities

import (
	"capital-gains-cli/internal/pkg/infrastructure/serialization"
	"capital-gains-cli/internal/pkg/utils"
	"math"
)

const (
	OperationTypeBuy           OperationType = "buy"
	OperationTypeSell          OperationType = "sell"
	MinOperationProfit         float64       = 0.00
	MinTaxFee                  float64       = 0.00
	QuantityMinStockQuantities float64       = 0.00
	QuantityMinCumulativeLoss  float64       = 0.00
)

type (
	OperationType string

	Operation struct {
		UnitCost      float64       `json:"unit-cost"`
		StockQuantity float64       `json:"quantity"`
		Operation     OperationType `json:"operation"`
	}

	StdinOperations []Operation

	TaxResult struct {
		Tax serialization.KeepZero `json:"tax"`
	}

	StoutOperations []TaxResult

	TaxState struct {
		WeightedMean    float64
		CumulativeLoss  float64
		StockQuantities float64
	}
)

func (taxState *TaxState) CalculateWeightedMean(operation Operation) {
	value := ((taxState.StockQuantities * taxState.WeightedMean) + (operation.StockQuantity * operation.UnitCost)) / (taxState.StockQuantities + operation.StockQuantity)

	taxState.WeightedMean = utils.RoundToTwoDecimalPlaces(value)
}

func (taxState *TaxState) AddStockQuantities(quantity float64) {
	taxState.StockQuantities += quantity
}

func (taxState *TaxState) RemoveStockQuantites(quantity float64) {
	newStockQuantities := taxState.StockQuantities - quantity

	taxState.StockQuantities = utils.SetMaxValue(
		QuantityMinStockQuantities,
		newStockQuantities,
	)
}

func (taxState *TaxState) SetCumulativeLoss(operationProfit float64) {
	if operationProfit < MinOperationProfit {
		taxState.CumulativeLoss += math.Abs(operationProfit)
		return
	}

	newCumulativeLoss := taxState.CumulativeLoss - operationProfit

	taxState.CumulativeLoss = utils.SetMaxValue(QuantityMinCumulativeLoss, newCumulativeLoss)
}
