package tests_test

import (
	"capital-gains-cli/internal/pkg/infrastructure/serialization"
	"encoding/json"
	"testing"
)

func TestJsonAdapters(t *testing.T) {
	testes := []struct {
		nome     string
		entrada  serialization.KeepZero
		esperado string
	}{
		{
			nome:     "integer number",
			entrada:  serialization.KeepZero(5),
			esperado: "5.00",
		},
		{
			nome:     "decimal number",
			entrada:  serialization.KeepZero(5.67),
			esperado: "5.67",
		},
		{
			nome:     "zero",
			entrada:  serialization.KeepZero(0),
			esperado: "0.00",
		},
		{
			nome:     "negative integer number",
			entrada:  serialization.KeepZero(-3),
			esperado: "-3.00",
		},
		{
			nome:     "negative decimal number",
			entrada:  serialization.KeepZero(-3.14),
			esperado: "-3.14",
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado, err := json.Marshal(tt.entrada)
			if err != nil {
				t.Errorf("Unexpected error when marshalling: %v", err)
			}

			if string(resultado) != tt.esperado {
				t.Errorf("Result = %v, expected = %v", string(resultado), tt.esperado)
			}
		})
	}
}
