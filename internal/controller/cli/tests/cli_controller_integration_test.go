package tests_test

import (
	"capital-gains-cli/internal/controller/cli"
	"capital-gains-cli/internal/domain/entities"
	"capital-gains-cli/internal/domain/strategies"
	usecase "capital-gains-cli/internal/useCase"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	TestFlags    = make(map[string]bool)
	buyStrategy  = strategies.NewBuyStrategy()
	sellStrategy = strategies.NewSellStrategy(0.20, 20000)

	mappedStrategies = usecase.MapStrategies{
		entities.OperationTypeBuy:  buyStrategy,
		entities.OperationTypeSell: sellStrategy,
	}

	operationHandler = usecase.NewOperationHandler(mappedStrategies)
	cliController    = cli.NewCliController(operationHandler)
)

func init() {
	if !TestFlags["operations"] {
		flag.String("operations", "", "Operations to be processed")
		TestFlags["operations"] = true
	}
}

func setupTestFlags(input string) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-operations", input}
}

func setupFileTestFlags(filepath string) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", filepath}
}

func TestCliControllerIntegration(t *testing.T) {
	t.Parallel()
	t.Run("(CASE#1) should not pay taxes when selling all shares with profit below 20000", func(t *testing.T) {

		input := `
			[
				{"operation":"buy", "unit-cost":10.00, "quantity": 100},
				{"operation":"sell", "unit-cost":20.00, "quantity": 50},
				{"operation":"sell", "unit-cost":5.00, "quantity": 50}
			]
		`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":0.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#2) should pay taxes only on profitable sale when having multiple operations", func(t *testing.T) {

		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":5.00, "quantity": 5000}
		]`

		expected := `[{"tax":0.00},{"tax":10000.00},{"tax":0.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#3) should calculate taxes when selling with loss and profit in different operations", func(t *testing.T) {

		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":5.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 3000}
		]`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":1000.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#4) should not pay taxes when selling with weighted average loss", func(t *testing.T) {
		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"buy", "unit-cost":25.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":15.00, "quantity": 10000}
		]`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":0.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#5) should pay taxes only on second sale when selling remaining shares with profit", func(t *testing.T) {

		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"buy", "unit-cost":25.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":15.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":25.00, "quantity": 5000}
		]`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":0.00},{"tax":10000.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#6) should use accumulated loss to offset future profits before charging taxes", func(t *testing.T) {

		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":2.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 2000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 2000},
			{"operation":"sell", "unit-cost":25.00, "quantity": 1000}
		]`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":0.00},{"tax":0.00},{"tax":3000.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#7) should use accumulated loss and calculate new profits in subsequent operations", func(t *testing.T) {

		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":2.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 2000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 2000},
			{"operation":"sell", "unit-cost":25.00, "quantity": 1000},
			{"operation":"buy", "unit-cost":20.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":15.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":30.00, "quantity": 4350},
			{"operation":"sell", "unit-cost":30.00, "quantity": 650}
		]`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":0.00},{"tax":0.00},{"tax":3000.00},{"tax":0.00},{"tax":0.00},{"tax":3700.00},{"tax":0.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#8) should calculate taxes on high profit sales in consecutive operations", func(t *testing.T) {

		input := `[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":50.00, "quantity": 10000},
			{"operation":"buy", "unit-cost":20.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":50.00, "quantity": 10000}
		]`

		expected := `[{"tax":0.00},{"tax":80000.00},{"tax":0.00},{"tax":60000.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, *response)
	})

	t.Run("(CASE#1+CASE#2) should handle both small and large operations with mixed tax scenarios", func(t *testing.T) {

		input := `
		[
			{"operation":"buy", "unit-cost":10.00, "quantity": 100},
			{"operation":"sell", "unit-cost":15.00, "quantity": 50},
			{"operation":"sell", "unit-cost":15.00, "quantity": 50}
		]
		[
			{"operation":"buy", "unit-cost":10.00, "quantity": 10000},
			{"operation":"sell", "unit-cost":20.00, "quantity": 5000},
			{"operation":"sell", "unit-cost":5.00, "quantity": 5000}
		]`

		expected := `[{"tax":0.00},{"tax":0.00},{"tax":0.00}]\n[{"tax":0.00},{"tax":10000.00},{"tax":0.00}]`

		setupTestFlags(input)

		response := cliController.HandlerStocks()

		assert.Equal(t, expected, strings.Replace(*response, "\n", "\\n", -1))
	})
}

func TestCliControllerIntegrationWithFile(t *testing.T) {
	t.Run("Should run multiples operations from a file", func(t *testing.T) {
		setupFileTestFlags("stubs/all_operations.txt")

		expectedBytes, err := os.ReadFile("stubs/expected_operation_responses.txt")
		if err != nil {
			t.Fatalf("Failed to read expected response file: %v", err)
		}
		expected := strings.TrimSpace(string(expectedBytes))
		expected = strings.ReplaceAll(expected, " ", "")
		expected = strings.ReplaceAll(expected, "\n", "\\n")
		expected = strings.ReplaceAll(expected, ",\\n{", ",{")
		expected = strings.ReplaceAll(expected, "},\\n", "},{")

		response := cliController.HandlerStocks()
		actual := strings.ReplaceAll(*response, "\n", "\\n")

		assert.Equal(t, expected, actual)
	})
}
