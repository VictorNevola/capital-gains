# Capital Gains CLI Calculator

A command-line tool for calculating capital gains taxes on stock operations.

## Requirements

- Go 1.23.4 or higher

## Dependencies

- github.com/golang/mock v1.6.0
- github.com/stretchr/testify v1.10.0

## Project Setup

1. Install dependencies:
```sh
make install
```

2. Run tests:
```sh
make test
```

3. Run tests with coverage:
```sh
make test-cover
```

4. Build all binaries:
```sh
make build-all
```
This will create executables in the `bin` directory:
- `/bin/linux/capital-gains`
- `/bin/windows/capital-gains.exe`
- `/bin/mac/capital-gains-intel`
- `/bin/mac/capital-gains-apple`


## Project Architecture

### Directory Structure

```
.
├── cmd/                    # Application entry point
├── internal/              
│   ├── controller/        # HTTP/CLI controllers
│   ├── domain/            # Business domain
│   │   ├── entities/      # Domain entities
│   │   ├── interfaces/    # Domain interfaces
│   │   └── strategies/    # Tax calculation strategies
│   ├── pkg/               # Shared packages
│   └── useCase/           # Application use cases
```

### Design Patterns

- **Strategy Pattern**: Used for different tax calculation methods (buy/sell)
- **Clean Architecture**: Separation of concerns with domain, use cases, and controllers
- **Dependency Injection**: For flexible component composition

### Key Components

1. **Tax Calculation Strategies**
   - `BuyStrategy`: Handles buy operations
   - `SellStrategy`: Handles sell operations with tax calculations

2. **Operation Handler**
   - Processes stock operations
   - Applies appropriate tax calculation strategy
   - Maintains operation state

3. **CLI Controller**
   - Handles command-line input/output
   - Processes JSON operation data
   - Returns tax calculation results


## Usage Example

### Running from source
```sh
go run cmd/main.go -operations='[{"operation":"buy", "unit-cost":10.00, "quantity":10000}, {"operation":"sell", "unit-cost":50.00, "quantity":10000}, {"operation":"buy", "unit-cost":20.00, "quantity":10000}, {"operation":"sell", "unit-cost":50.00, "quantity":10000}]'
```
```sh
go run cmd/main.go 'operations.txt'
```
### Running from built binary
### Running on Linux
```sh
./bin/linux/capital-gains-linux-amd64 -operations='[{"operation":"buy", "unit-cost":10.00, "quantity":10000}, {"operation":"sell", "unit-cost":50.00, "quantity":10000}, {"operation":"buy", "unit-cost":20.00, "quantity":10000}, {"operation":"sell", "unit-cost":50.00, "quantity":10000}]'
```
```sh
./bin/linux/capital-gains-linux-amd64 'test-operations.txt'
```

### Running on Windows
**Note: It's recommended to use CMD instead of PowerShell for better compatibility with command-line arguments.**

Using CMD:
```cmd
# Navigate to the binary directory
cd bin/windows

# Using command-line arguments (recommended to use double quotes and escape characters)
capital-gains.exe -operations="[{\"operation\":\"buy\",\"unit-cost\":10.00,\"quantity\":10000},{\"operation\":\"sell\",\"unit-cost\":50.00,\"quantity\":10000}]"

# Alternative: Using an input file (simpler approach)
capital-gains.exe "../../test-operations.txt"
```

### Running on macOS
For Intel Mac:
```sh
./bin/mac/capital-gains-intel -operations='[{"operation":"buy", "unit-cost":10.00, "quantity":10000}, {"operation":"sell", "unit-cost":50.00, "quantity":10000}]'

./bin/mac/capital-gains-intel 'test-operations.txt'
```

For Apple Silicon (M1/M2):
```sh
./bin/mac/capital-gains-apple -operations='[{"operation":"buy", "unit-cost":10.00, "quantity":10000}, {"operation":"sell", "unit-cost":50.00, "quantity":10000}]'

./bin/mac/capital-gains-apple 'test-operations.txt'
```
