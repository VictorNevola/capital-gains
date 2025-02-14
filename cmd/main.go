package main

import (
	cli "capital-gains-cli/internal/controller/cli"
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/domain/strategies"
	usecase "capital-gains-cli/internal/useCase"
	"fmt"
)

const (
	feeTax         = 0.20
	minAmountToTax = 20000
)

func main() {
	buyStrategy := strategies.NewBuyStrategy()
	sellStrategy := strategies.NewSellStrategy(
		feeTax,
		minAmountToTax,
	)

	operationHandler := usecase.NewOperationHandler(usecase.MapStrategies{
		entities.OperationTypeBuy:  buyStrategy,
		entities.OperationTypeSell: sellStrategy,
	})

	cliController := cli.NewCliController(operationHandler)

	result := cliController.HandlerStocks()

	fmt.Println(*result)
}
