package tests_test

import (
	"capital-gains-cli/internal/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundToTwoDecimalPlaces(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{input: 1.2345, expected: 1.23},
		{input: 1.2355, expected: 1.24},
		{input: 1.2, expected: 1.20},
		{input: -1.2345, expected: -1.23},
		{input: -1.2355, expected: -1.24},
	}

	for _, test := range tests {
		result := utils.RoundToTwoDecimalPlaces(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestSetMaxValue(t *testing.T) {
	tests := []struct {
		min      float64
		max      float64
		expected float64
	}{
		{min: 0, max: 5, expected: 5},
		{min: 10, max: 5, expected: 10},
		{min: -5, max: -10, expected: -5},
		{min: -10, max: -5, expected: -5},
	}

	for _, test := range tests {
		result := utils.SetMaxValue(test.min, test.max)
		assert.Equal(t, test.expected, result)
	}
}
