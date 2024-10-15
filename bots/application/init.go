package application

import (
	"context"
	"fmt"
	"time"

	"github.com/juankohler/crypto-bot/bots/domain"
	"github.com/juankohler/crypto-bot/libs/go/errors"
	"github.com/juankohler/crypto-bot/libs/go/logs"
	"github.com/shopspring/decimal"
)

type InitInput struct {
	// Name        string `json:"name"`
}

type Init struct {
	providerRepository domain.ProviderRepository
}

func NewInit(
	providerRepository domain.ProviderRepository,
) *Init {
	return &Init{
		providerRepository: providerRepository,
	}
}

func (s *Init) Exec(ctx context.Context, input *InitInput) error {
	name := "JUANCHO"
	currency := domain.CurrencyUSDT
	targetCurrency := domain.CurrencyBTC
	takeProfitPercentaje := decimal.NewFromFloat(0.005)
	initialCapital := decimal.NewFromFloat(1000)
	delta := decimal.NewFromFloat(200)
	monitorInterval := 20 * time.Second

	juancho, err := domain.CreateBot(name, currency, targetCurrency, takeProfitPercentaje, initialCapital, delta, monitorInterval)
	if err != nil {
		return err
	}

	logs.Info(ctx, fmt.Sprintf("Bot %s inicializado, capital_inicial: %s, delta: %s, take_profit_percentage: %s", juancho.Name, juancho.InitialCapital.String(), juancho.Delta.String(), juancho.TakeProfitPercentaje.String()))

	name2 := "ALE"
	currency2 := domain.CurrencyUSDT
	targetCurrency2 := domain.CurrencyBTC
	takeProfitPercentaje2 := decimal.NewFromFloat(0.005)
	initialCapital2 := decimal.NewFromFloat(1000)
	delta2 := decimal.NewFromFloat(200)
	monitorInterval2 := 20 * time.Second

	ale, err := domain.CreateBot(name2, currency2, targetCurrency2, takeProfitPercentaje2, initialCapital2, delta2, monitorInterval2)
	if err != nil {
		return err
	}

	logs.Info(ctx, fmt.Sprintf("Bot %s inicializado, capital_inicial: %s, delta: %s, take_profit_percentage: %s", ale.Name, ale.InitialCapital.String(), ale.Delta.String(), ale.TakeProfitPercentaje.String()))

	go s.executeBots(ctx, juancho, ale)

	return nil
}

func (s *Init) executeBots(ctx context.Context, juancho *domain.Bot, ale *domain.Bot) {
	first := true
	for {
		if !first {
			time.Sleep(juancho.MonitorInterval)
		}

		currentPrice, err := s.providerRepository.GetPrice(ctx, juancho.TargetCurrency, juancho.Currency)
		if err != nil {
			logs.Error(ctx, "could not get price", logs.NewAttr("error", err))
			continue
		}

		// logs.Info(ctx, fmt.Sprintf("Precio BTC/USDT: %s", currentPrice.Price.String()))

		err = s.executeJuanchoStrategy(ctx, juancho, currentPrice.Price)
		if err != nil {
			logs.Error(ctx, "error in JUANCHO strategy", logs.NewAttr("error", err))
		}

		err = s.executeAleStrategy(ctx, ale, currentPrice.Price)
		if err != nil {
			logs.Error(ctx, "error in ALE strategy", logs.NewAttr("error", err))
		}

		first = false
	}
}

func (s *Init) executeJuanchoStrategy(ctx context.Context, bot *domain.Bot, currentPrice decimal.Decimal) error {
	/** TODO: only for simulate sell with test */
	bot.RemoveOrdersBelowPrice(ctx, currentPrice)

	priceRange := bot.CalculatePriceRange(currentPrice)
	if bot.HasOpenOrderWithPriceRange(priceRange) {
		return nil
	}

	quoteAmount := bot.InitialCapital.Div(decimal.NewFromFloat(50))
	newOrder, err := bot.GenerateOrder(currentPrice, priceRange, quoteAmount)
	if err != nil {
		return errors.Wrap(domain.ErrInternal, err, "could not generate order")
	}

	externalId, err := s.providerRepository.CreateOrderInProvider(ctx, newOrder, bot.Name)
	if err != nil {
		return errors.Wrap(domain.ErrInternal, err, "could not create order in provider")
	}

	newOrder.AddExternalId(externalId)

	return nil
}

func (s *Init) executeAleStrategy(ctx context.Context, bot *domain.Bot, currentPrice decimal.Decimal) error {
	/** TODO: only for simulate sell with test */
	bot.RemoveOrdersBelowPrice(ctx, currentPrice)

	first := true
	if bot.LastSalePrice != nil {
		first = false
	}

	priceHasDroppedBelowDelta := false
	if !first && bot.LastSalePrice.GreaterThan(currentPrice) {
		priceDifference := bot.LastSalePrice.Sub(currentPrice)
		priceHasDroppedBelowDelta = priceDifference.GreaterThan(bot.Delta)
	}

	if bot.HasOpenOrder() || (!first && !priceHasDroppedBelowDelta) {
		return nil
	}

	priceRange := bot.CalculatePriceRange(currentPrice)
	quoteAmount := bot.AvailableCapital
	newOrder, err := bot.GenerateOrder(currentPrice, priceRange, quoteAmount)
	if err != nil {
		return errors.Wrap(domain.ErrInternal, err, "could not generate order")
	}

	externalId, err := s.providerRepository.CreateOrderInProvider(ctx, newOrder, bot.Name)
	if err != nil {
		return errors.Wrap(domain.ErrInternal, err, "could not create order in provider")
	}

	newOrder.AddExternalId(externalId)

	return nil
}
