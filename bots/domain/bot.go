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
	ID                   models.ID
	Name                 string
	TakeProfitPercentaje decimal.Decimal
	InitialCapital       decimal.Decimal
	AvailableCapital     decimal.Decimal
	InvestedCapital      decimal.Decimal
	TotalCapital         decimal.Decimal
	Currency             string
	TargetCurrency       string
	Delta                decimal.Decimal
	MonitorInterval      time.Duration
	Timestamps           models.Timestamps
	Version              models.Version
	OpenOrders           []*Order
	LastSalePrice        *decimal.Decimal

	mu sync.Mutex
}

func NewBot(
	id models.ID,
	name string,
	currency string,
	targetCurrency string,
	takeProfitPercentaje decimal.Decimal,
	initialCapital decimal.Decimal,
	availableCapital decimal.Decimal,
	investedCapital decimal.Decimal,
	totalCapital decimal.Decimal,
	delta decimal.Decimal,
	monitorInterval time.Duration,
	openOrders []*Order,
	lastSalePrice *decimal.Decimal,
	timestamps models.Timestamps,
	version models.Version,
) (*Bot, error) {
	entity := &Bot{
		ID:                   id,
		Name:                 name,
		TakeProfitPercentaje: takeProfitPercentaje,
		InitialCapital:       initialCapital,
		AvailableCapital:     availableCapital,
		InvestedCapital:      investedCapital,
		TotalCapital:         totalCapital,
		Currency:             currency,
		TargetCurrency:       targetCurrency,
		Delta:                delta,
		MonitorInterval:      monitorInterval,
		OpenOrders:           openOrders,
		LastSalePrice:        lastSalePrice,
		Timestamps:           timestamps,
		Version:              version,
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
	takeProfitPercentaje decimal.Decimal,
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
	var lastSalePrice *decimal.Decimal

	entity, err := NewBot(
		id,
		name,
		currency,
		targetCurrency,
		takeProfitPercentaje,
		initialCapital,
		availableCapital,
		investedCapital,
		totalCapital,
		delta,
		monitorInterval,
		openOrders,
		lastSalePrice,
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

func (s *Bot) HasOpenOrderWithPriceRange(priceRange int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, order := range s.OpenOrders {
		if order.PriceRange == priceRange {
			return true
		}
	}
	return false
}

func (s *Bot) HasOpenOrder() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.OpenOrders) > 0
}

func (s *Bot) GenerateOrder(currentPrice decimal.Decimal, priceRange int, initialQuoteAmount decimal.Decimal) (*Order, error) {
	quantity := initialQuoteAmount.Div(currentPrice)
	takeProfit := currentPrice.Add(currentPrice.Mul(s.TakeProfitPercentaje))
	finalQuoteAmount := quantity.Mul(takeProfit)
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
		initialQuoteAmount,
		finalQuoteAmount,
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
	s.InvestedCapital = s.InvestedCapital.Add(newOrder.InitialQuoteAmount)
	s.AvailableCapital = s.AvailableCapital.Sub(newOrder.InitialQuoteAmount)
	s.TotalCapital = s.InvestedCapital.Add(s.AvailableCapital)
	s.OpenOrders = append(s.OpenOrders, newOrder)
	s.mu.Unlock()

	return newOrder, nil
}

func (s *Bot) RemoveOrdersBelowPrice(ctx context.Context, price decimal.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var filteredOrders []*Order
	for _, order := range s.OpenOrders {
		if order.TakeProfitPrice.GreaterThan(price) {
			filteredOrders = append(filteredOrders, order)
		} else {
			logs.Info(ctx, fmt.Sprintf("%s: Venta %s BTC a %s USDT (%s USDT)", s.Name, order.Quantity.String(), order.TakeProfitPrice.String(), order.FinalQuoteAmount.String()))
			s.InvestedCapital = s.InvestedCapital.Sub(order.InitialQuoteAmount)
			s.AvailableCapital = s.AvailableCapital.Add(order.FinalQuoteAmount)
			s.TotalCapital = s.InvestedCapital.Add(s.AvailableCapital)
			lastSalePrice := price
			s.LastSalePrice = &lastSalePrice

			logs.Info(ctx, fmt.Sprintf("%s: SALDO TOTAL: %s", s.Name, s.TotalCapital.String()))
		}
	}

	s.OpenOrders = filteredOrders
}
