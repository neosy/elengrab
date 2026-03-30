package reflection

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructFieldName(
	t *testing.T) {
	tests := []struct {
		name          string
		structure     any
		field         any
		tag           string
		expected      string
		expectedError error
	}{
		{
			"Getting field name trade_id",
			&tradeEmpty,
			&tradeEmpty.TradeId,
			"db",
			"trade_id",
			nil,
		},
		{
			"Getting field name instrumentId",
			&tradeEmpty,
			&tradeEmpty.InstrumentId,
			"json",
			"instrumentId",
			nil,
		},
		{
			"Getting field name qty",
			&tradeEmpty,
			&tradeEmpty.Qty,
			"db",
			"qty",
			nil,
		},
		{
			"Incorrect tag provided",
			&tradeEmpty,
			&tradeEmpty.TradeId,
			"dbb",
			"",
			ErrStructureFieldNameEmpty,
		},
		{
			"First parameter is not a pointer",
			tradeEmpty,
			&tradeEmpty.TradeId,
			"db",
			"trade_id",
			ErrFirstArgumentTypeMustPointerStructure,
		},
		{
			"Second parameter is not a pointer",
			&tradeEmpty,
			tradeEmpty.TradeId,
			"db",
			"trade_id",
			ErrSecondArgumentTypeMustPointerFieldStructure,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result, err := StructFieldName(tt.structure, tt.field, tt.tag)

			if tt.expectedError != nil {
				if err != nil {
					assert.EqualError(t, err, tt.expectedError.Error(), fmt.Sprintf("Test %q: expected error", tt.name))
				} else {
					assert.Fail(t, fmt.Sprintf("Test %q: expected error %v, but it did not occur", tt.name, tt.expectedError))
				}
				return
			}

			if err != nil {
				t.Errorf("Test %q: unexpected error, got %v", tt.name, err)
			} else {
				assert.Equal(t, tt.expected, result, fmt.Sprintf("Test %q: expected %s, got %s", tt.name, tt.expected, result))
			}
		})
	}
}
