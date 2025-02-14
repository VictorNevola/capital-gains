package cli

import (
	"capital-gains-cli/internal/domain/entities"
	usecase "capital-gains-cli/internal/useCase"
	"encoding/json"
	"flag"
	"os"
	"strings"
)

type (
	CliController struct {
		stocksOperationHandler *usecase.OperationHandler
	}
)

func NewCliController(stocksOperationHandler *usecase.OperationHandler) *CliController {
	return &CliController{
		stocksOperationHandler: stocksOperationHandler,
	}
}

func (c *CliController) HandlerStocks() *string {
	operations := c.getInput()
	stdinOperations, err := c.inputReader(operations)
	if err != nil {
		panic(err)
	}

	allTaxResults, err := c.stocksOperationHandler.Handler(stdinOperations)
	if err != nil {
		panic(err)
	}

	response, err := c.outputWriter(allTaxResults)
	if err != nil {
		panic(err)
	}

	return response
}

func (c *CliController) getInput() string {
	operations := c.getFromFlag()
	operations += c.getFromFile()

	return c.sanitizeOperations(operations)
}

func (c CliController) getFromFlag() string {
	inputFlag := flag.String("operations", "", "Operations to be processed")
	flag.Parse()

	return *inputFlag
}

func (c CliController) getFromFile() string {
	if len(os.Args) > 1 && strings.HasSuffix(os.Args[1], ".txt") {
		content, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		return string(content)
	}

	return ""
}

func (c CliController) sanitizeOperations(operations string) string {
	input := strings.TrimSpace(operations)

	quotes := []string{"'", "\""}
	for _, quote := range quotes {
		input = strings.Trim(input, quote)
	}

	return input
}

func (h *CliController) inputReader(input string) ([][]entities.Operation, error) {
	operationsArrays := strings.SplitAfter(input, "]")
	var operationsGroups [][]entities.Operation

	for _, arr := range operationsArrays {
		if arr == "" {
			continue
		}

		var operations entities.StdinOperations
		error := json.Unmarshal([]byte(arr), &operations)
		if error != nil {
			return nil, error
		}

		operationsGroups = append(operationsGroups, operations)
	}

	return operationsGroups, nil
}

func (h *CliController) outputWriter(allTaxResults [][]entities.TaxResult) (*string, error) {
	var formattedArrays []string

	for _, taxResults := range allTaxResults {
		jsonResult, err := json.Marshal(taxResults)
		if err != nil {
			return nil, err
		}

		formattedArrays = append(formattedArrays, string(jsonResult))
	}

	result := strings.Join(formattedArrays, "\n")
	return &result, nil
}
