package usecase

import (
	"capital-gains-cli/internal/domain/entities"
	interfaces "capital-gains-cli/internal/domain/interfaces"
)

type (
	MapStrategies map[entities.OperationType]interfaces.TaxCalculator

	OperationHandler struct {
		strategies MapStrategies
	}
)

func NewOperationHandler(mappedStrategies MapStrategies) *OperationHandler {
	return &OperationHandler{
		strategies: mappedStrategies,
	}
}

func (h *OperationHandler) Handler(operations [][]entities.Operation) ([][]entities.TaxResult, error) {
	allTaxResults := make([][]entities.TaxResult, len(operations))

	for index, operation := range operations {
		taxsResults, err := h.processOperation(operation)
		if err != nil {
			return nil, err
		}
		allTaxResults[index] = taxsResults
	}

	return allTaxResults, nil
}

func (h *OperationHandler) processOperation(operations entities.StdinOperations) ([]entities.TaxResult, error) {
	taxsResults := make([]entities.TaxResult, len(operations))
	operationsState := entities.TaxState{}

	for index, operation := range operations {
		strategy, ok := h.strategies[operation.Operation]
		if !ok {
			return nil, ErrOperationTypeNotSupported
		}

		taxsResults[index] = strategy.CalculateTax(operation, &operationsState)
	}

	return taxsResults, nil
}
