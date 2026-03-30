package reflection

import "github.com/google/uuid"

type tradeEntity struct {
	TradeId      uuid.UUID `db:"trade_id" json:"tradeId"`
	InstrumentId uuid.UUID `db:"instrument_id" json:"instrumentId"`
	Price        uint64    `db:"price" json:"price"`
	Qty          uint64    `db:"qty" json:"qty"`
}

var tradeEmpty = tradeEntity{}
