package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFixedUp(t *testing.T) {

	input := 199.546
	expected := 199.55

	result := ToFixed(input)

	assert.Equal(t, expected, result)
}

func TestToFixedDown(t *testing.T) {

	input := 199.544
	expected := 199.54

	result := ToFixed(input)

	assert.Equal(t, expected, result)
}

func TestToFixedEqual(t *testing.T) {

	input := 199.54
	expected := 199.54

	result := ToFixed(input)

	assert.Equal(t, expected, result)
}

func TestToFixedNegative(t *testing.T) {

	input := -199.544
	expected := -199.54

	result := ToFixed(input)

	assert.Equal(t, expected, result)
}

func TestToFixedZero(t *testing.T) {

	input := 0.00
	expected := 0.00

	result := ToFixed(input)

	assert.Equal(t, expected, result)
}
