package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToKopecks(t *testing.T) {
	tests := []struct {
		value  float64
		wanted int64
	}{
		{
			value:  5.32,
			wanted: 532,
		},
		{
			value:  1.0,
			wanted: 100,
		},
		{
			value:  1.1,
			wanted: 110,
		},
		{
			value:  1.119,
			wanted: 111,
		},
	}

	for i, test := range tests {
		test := test

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			res := ConvertFloatToKopecks(test.value)

			assert.Equal(t, test.wanted, res)
		})
	}
}

func TestConvertKopecksToFloat(t *testing.T) {
	tests := []struct {
		value  int64
		wanted float64
	}{
		{
			value:  532,
			wanted: 5.32,
		},
		{
			value:  100,
			wanted: 1.00,
		},
		{
			value:  110,
			wanted: 1.1,
		},
		{
			value:  111,
			wanted: 1.11,
		},
	}

	for i, test := range tests {
		test := test

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			res := ConvertKopecksToFloat(test.value)

			assert.Equal(t, test.wanted, res)
		})
	}
}
