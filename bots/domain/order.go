package domain

import (
	"github.com/juankohler/crypto-bot/libs/go/models"
	"github.com/shopspring/decimal"
)

type OrderRepository interface {
	// FindByID(ctx context.Context, id models.ID) (*Order, error)
	// Save(ctx context.Context, order *Order) error
}

const (
	OrderStatusPending   = "PENDING"
	OrderStatusOpen      = "OPEN"
	OrderStatusCompleted = "COMPLETED"
)

type Order struct {
	ID                 models.ID
	BotID              models.ID
	Symbol             string
	Quantity           decimal.Decimal
	InitialQuoteAmount decimal.Decimal
	FinalQuoteAmount   decimal.Decimal
	EntryPrice         decimal.Decimal
	TakeProfitPrice    decimal.Decimal
	ExternalId         *string
	Status             string
	PriceRange         int
	Timestamps         models.Timestamps
	Version            models.Version
}

func NewOrder(
	id models.ID,
	botID models.ID,
	symbol string,
	quantity decimal.Decimal,
	initialQuoteAmount decimal.Decimal,
	finalQuoteAmount decimal.Decimal,
	entryPrice decimal.Decimal,
	takeProfitPrice decimal.Decimal,
	externalId *string,
	status string,
	priceRange int,
	timestamps models.Timestamps,
	version models.Version,
) (*Order, error) {
	entity := &Order{
		ID:                 id,
		BotID:              botID,
		Symbol:             symbol,
		Quantity:           quantity,
		InitialQuoteAmount: initialQuoteAmount,
		FinalQuoteAmount:   finalQuoteAmount,
		EntryPrice:         entryPrice,
		TakeProfitPrice:    takeProfitPrice,
		ExternalId:         externalId,
		Status:             status,
		PriceRange:         priceRange,
		Timestamps:         timestamps,
		Version:            version,
	}

	return entity, nil
}

// func (s *Order) updated() {
// 	s.Timestamps = s.Timestamps.Update()
// 	s.Version = s.Version.Update()
// }

func (s *Order) AddExternalId(externalId string) {
	s.ExternalId = &externalId
	s.Status = OrderStatusOpen
}
