package reflection

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type structTest1 struct {
	TradeId      string `db:"trade_id" json:"tradeId"`
	InstrumentId string `db:"instrument_id" json:"instrumentId"`
	Price        uint64 `db:"price" json:"price"`
	Qty          uint64 `db:"qty" json:"qty"`
	Amount       string `db:"amount" json:"amount"`
}

type structTest2 struct {
	TradeId            string     `db:"trade_id" json:"tradeId"`
	InstrumentId       string     `db:"instrument_id" json:"instrumentId"`
	Price              uint64     `db:"price" json:"price"`
	Qty                uint64     `db:"qty" json:"qty"`
	Amount             string     `db:"amount" json:"amount"`
	GateId             *uuid.UUID `db:"gate_id"`
	GateOrderId        *string    `db:"gate_order_id"`
	GatePaymentDetails string     `db:"gate_payment_details"`
	SettlementType     string     `db:"settlement_type"`
	Rate               string     `db:"rate"`
	RateTimestamp      time.Time  `db:"rate_timestamp"`
	CreatedBy          uuid.UUID  `db:"created_by"`
	CreatedAt          time.Time  `db:"created_at" insert:"false"`
	UpdatedAt          time.Time  `db:"updated_at" insert:"false"`
	CancelledAt        *time.Time `db:"canceled_at"`
}

func runTestFirstOf5() {
	var sTest = structTest1{}

	StructFieldName(&sTest, &sTest.TradeId, "db")
}

func runTestFirstOf15() {
	var sTest = structTest2{}

	StructFieldName(&sTest, &sTest.TradeId, "db")
}

func runTestFifthOf5() {
	var sTest = structTest1{}

	StructFieldName(&sTest, &sTest.Amount, "db")
}

func runTestFifthOf15() {
	var sTest = structTest2{}

	StructFieldName(&sTest, &sTest.Amount, "db")
}

func runTestLastOf5() {
	var sTest = structTest1{}

	StructFieldName(&sTest, &sTest.Amount, "db")
}

func runTestLastOf15() {
	var sTest = structTest2{}

	StructFieldName(&sTest, &sTest.CancelledAt, "db")
}

func runTest5of5() {
	var sTest = structTest1{}

	StructFieldName(&sTest, &sTest.TradeId, "db")
	StructFieldName(&sTest, &sTest.InstrumentId, "db")
	StructFieldName(&sTest, &sTest.Price, "db")
	StructFieldName(&sTest, &sTest.Qty, "db")
	StructFieldName(&sTest, &sTest.Amount, "db")
}

func runTest15of15() {
	var sTest = structTest2{}

	StructFieldName(&sTest, &sTest.TradeId, "db")
	StructFieldName(&sTest, &sTest.InstrumentId, "db")
	StructFieldName(&sTest, &sTest.Price, "db")
	StructFieldName(&sTest, &sTest.Qty, "db")
	StructFieldName(&sTest, &sTest.Amount, "db")
	StructFieldName(&sTest, &sTest.GateId, "db")
	StructFieldName(&sTest, &sTest.GateOrderId, "db")
	StructFieldName(&sTest, &sTest.GatePaymentDetails, "db")
	StructFieldName(&sTest, &sTest.SettlementType, "db")
	StructFieldName(&sTest, &sTest.Rate, "db")
	StructFieldName(&sTest, &sTest.RateTimestamp, "db")
	StructFieldName(&sTest, &sTest.CreatedBy, "db")
	StructFieldName(&sTest, &sTest.CreatedAt, "db")
	StructFieldName(&sTest, &sTest.UpdatedAt, "db")
	StructFieldName(&sTest, &sTest.CancelledAt, "db")
}

func BenchmarkFields5of5(b *testing.B) {
	for b.Loop() {
		runTest5of5()
	}
}

func BenchmarkFields15of15(b *testing.B) {
	for b.Loop() {
		runTest15of15()
	}
}

func BenchmarkFieldFirstOf5(b *testing.B) {
	for b.Loop() {
		runTestFirstOf5()
	}
}

func BenchmarkFieldFirstOf15(b *testing.B) {
	for b.Loop() {
		runTestFirstOf15()
	}
}

func BenchmarkFieldsFifthOf5(b *testing.B) {
	for b.Loop() {
		runTestFifthOf5()
	}
}

func BenchmarkFieldsFifthOf15(b *testing.B) {
	for b.Loop() {
		runTestFifthOf15()
	}
}

func BenchmarkFieldsLastOf5(b *testing.B) {
	for b.Loop() {
		runTestLastOf5()
	}
}

func BenchmarkFieldsLastOf15(b *testing.B) {
	for b.Loop() {
		runTestLastOf15()
	}
}
