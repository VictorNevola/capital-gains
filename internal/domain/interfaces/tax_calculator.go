package domain

import "capital-gains-cli/internal/domain/entities"

type TaxCalculator interface {
	CalculateTax(operation entities.Operation, taxState *entities.TaxState) entities.TaxResult
}
