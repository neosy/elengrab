package reflection

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type tradeEntity struct {
	TradeId      uuid.UUID `db:"trade_id" json:"tradeId"`
	InstrumentId uuid.UUID `db:"instrument_id" json:"instrumentId"`
	Price        uint64    `db:"price" json:"price"`
	Qty          uint64    `db:"qty" json:"qty"`
}

var tradeEmpty = tradeEntity{}

func TestStructFieldName(
	t *testing.T) {
	tests := []struct {
		name          string
		structure     interface{}
		field         interface{}
		tag           string
		expected      string
		expectedError error
	}{
		{
			"Получение названия поля trade_id",
			&tradeEmpty,
			&tradeEmpty.TradeId,
			"db",
			"trade_id",
			nil,
		},
		{
			"Получение названия поля instrumentId",
			&tradeEmpty,
			&tradeEmpty.InstrumentId,
			"json",
			"instrumentId",
			nil,
		},
		{
			"Получение названия поля qty",
			&tradeEmpty,
			&tradeEmpty.Qty,
			"db",
			"qty",
			nil,
		},
		{
			"Не верно передан tag",
			&tradeEmpty,
			&tradeEmpty.TradeId,
			"db",
			"",
			ErrStructureFieldNameEmpty,
		},
		{
			"Первый параметр не указатель",
			tradeEmpty,
			&tradeEmpty.TradeId,
			"db",
			"trade_id",
			ErrFirstArgumentTypeMustPointerStructure,
		},
		{
			"Второй параметр не указатель",
			&tradeEmpty,
			tradeEmpty.TradeId,
			"db",
			"trade_id",
			ErrSecondArgumentTypeMustPointerFieldStructure,
		},
	}

	for _, tt := range tests {
		result, err := StructFieldName(tt.structure, tt.field, tt.tag)

		if tt.expectedError != nil {
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error(), "Ожидаемая ошибка %v, получил: %v", tt.expectedError, err)
			}

			assert.Error(t, err, "Ожидалась ошибка %v, но её не было", tt.expectedError)

			continue
		}

		if err != nil {
			t.Errorf("Ошибка не ожидалась, получено %v", err)
		} else {
			assert.Equal(t, tt.expected, result, "Ожидал %d, получил %d", tt.expected, result)
		}
	}
}
