package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/juankohler/crypto-bot/libs/go/errors"
	"github.com/juankohler/crypto-bot/libs/go/logs"
	"github.com/juankohler/crypto-bot/libs/go/models"
	"github.com/shopspring/decimal"
)

type BotRepository interface {
	FindByID(ctx context.Context, id models.ID) (*Bot, error)
	Save(ctx context.Context, user *Bot) error
}

type Bot struct {
	ID               models.ID
	Name             string
	InitialCapital   decimal.Decimal
	AvailableCapital decimal.Decimal
	InvestedCapital  decimal.Decimal
	TotalCapital     decimal.Decimal
	Currency         string
	TargetCurrency   string
	Delta            decimal.Decimal
	MonitorInterval  time.Duration
	Timestamps       models.Timestamps
	Version          models.Version
	OpenOrders       []*Order

	mu sync.Mutex
}

func NewBot(
	id models.ID,
	name string,
	currency string,
	targetCurrency string,
	initialCapital decimal.Decimal,
	availableCapital decimal.Decimal,
	investedCapital decimal.Decimal,
	totalCapital decimal.Decimal,
	delta decimal.Decimal,
	monitorInterval time.Duration,
	openOrders []*Order,
	timestamps models.Timestamps,
	version models.Version,
) (*Bot, error) {
	entity := &Bot{
		ID:               id,
		Name:             name,
		InitialCapital:   initialCapital,
		AvailableCapital: availableCapital,
		InvestedCapital:  investedCapital,
		TotalCapital:     totalCapital,
		Currency:         currency,
		TargetCurrency:   targetCurrency,
		Delta:            delta,
		MonitorInterval:  monitorInterval,
		OpenOrders:       openOrders,
		Timestamps:       timestamps,
		Version:          version,
	}

	return entity, nil
}

// func (t *Bot) updated() {
// 	t.Timestamps = t.Timestamps.Update()
// 	t.Version = t.Version.Update()
// }

func CreateBot(
	name string,
	currency string,
	targetCurrency string,
	initialCapital decimal.Decimal,
	delta decimal.Decimal,
	monitorInterval time.Duration,
) (*Bot, error) {
	id, err := models.GenerateNanoID(10)
	if err != nil {
		errors.Wrap(ErrInternal, err, "could not generate nano id")
	}

	availableCapital := initialCapital
	investedCapital := decimal.NewFromFloat(0)
	totalCapital := availableCapital.Add(investedCapital)
	openOrders := []*Order{}

	entity, err := NewBot(
		id,
		name,
		targetCurrency,
		currency,
		initialCapital,
		availableCapital,
		investedCapital,
		totalCapital,
		delta,
		monitorInterval,
		openOrders,
		models.CreateTimestamps(),
		models.CreateVersion(),
	)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *Bot) CalculatePriceRange(currentPrice decimal.Decimal) int {
	return int(currentPrice.Div(s.Delta).Floor().IntPart())
}

func (s *Bot) HasOpenOrder(priceRange int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, order := range s.OpenOrders {
		if order.PriceRange == priceRange {
			return true
		}
	}
	return false
}

func (s *Bot) GenerateOrder(currentPrice decimal.Decimal, priceRange int) (*Order, error) {
	quoteAmount := s.InitialCapital.Div(decimal.NewFromFloat(100))
	quantity := quoteAmount.Div(currentPrice)
	takeProfit := currentPrice.Mul(decimal.NewFromFloat(1.02))
	status := OrderStatusPending
	symbol := s.TargetCurrency + "/" + s.Currency
	orderId, err := models.GenerateNanoID(14)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err, "could not generate order id")
	}

	newOrder, err := NewOrder(
		orderId,
		s.ID,
		symbol,
		quantity,
		quoteAmount,
		currentPrice,
		takeProfit,
		nil,
		status,
		priceRange,
		models.CreateTimestamps(),
		models.CreateVersion(),
	)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.InvestedCapital = s.InvestedCapital.Add(newOrder.QuoteAmount)
	s.AvailableCapital = s.AvailableCapital.Sub(newOrder.QuoteAmount)
	s.TotalCapital = s.InvestedCapital.Add(s.AvailableCapital)
	s.OpenOrders = append(s.OpenOrders, newOrder)
	s.mu.Unlock()

	return newOrder, nil
}

func (s *Bot) RemoveOrdersBelowPriceRange(ctx context.Context, priceRange int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var filteredOrders []*Order
	for _, order := range s.OpenOrders {
		if order.PriceRange >= priceRange {
			filteredOrders = append(filteredOrders, order)
		} else {
			logs.Info(ctx, "VENTA", logs.NewAttr("order", order))
			s.InvestedCapital = s.InvestedCapital.Sub(order.QuoteAmount)
			s.AvailableCapital = s.AvailableCapital.Add(order.QuoteAmount)
			s.TotalCapital = s.InvestedCapital.Add(s.AvailableCapital)

			logs.Info(ctx, fmt.Sprintf("SALDO TOTAL: %s", s.TotalCapital.String()))
		}
	}

	s.OpenOrders = filteredOrders
}
