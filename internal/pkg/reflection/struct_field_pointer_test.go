package reflection

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStructFieldPointer(t *testing.T) {
	// Prepare an empty tradeEntity
	trade := tradeEntity{}

	tests := []struct {
		name          string
		structure     any
		fieldName     string
		tag           string
		expectedType  any
		expectedError error
	}{
		{
			name:          "Get TradeId by field name",
			structure:     &trade,
			fieldName:     "TradeId",
			tag:           "",
			expectedType:  (*uuid.UUID)(nil),
			expectedError: nil,
		},
		{
			name:          "Get InstrumentId by json tag",
			structure:     &trade,
			fieldName:     "instrumentId",
			tag:           "json",
			expectedType:  (*uuid.UUID)(nil),
			expectedError: nil,
		},
		{
			name:          "Get Qty by db tag",
			structure:     &trade,
			fieldName:     "qty",
			tag:           "db",
			expectedType:  (*uint64)(nil),
			expectedError: nil,
		},
		{
			name:          "Field not found with wrong name",
			structure:     &trade,
			fieldName:     "UnknownField",
			tag:           "",
			expectedType:  nil,
			expectedError: ErrFieldNotFound,
		},
		{
			name:          "Field not found with wrong tag",
			structure:     &trade,
			fieldName:     "wrong_tag",
			tag:           "json",
			expectedType:  nil,
			expectedError: ErrFieldNotFound,
		},
		{
			name:          "Invalid input, not a pointer",
			structure:     trade,
			fieldName:     "TradeId",
			tag:           "",
			expectedType:  nil,
			expectedError: ErrInputMustBePointerToStruct,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result, err := StructFieldPointer(tt.structure, tt.fieldName, tt.tag)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Check type of pointer
			assert.IsType(t, tt.expectedType, result)
		})
	}
}
