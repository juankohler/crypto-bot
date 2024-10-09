package application

import (
	"context"
	"fmt"
	"time"

	"github.com/juankohler/crypto-bot/bots/domain"
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
	logs.Info(ctx, "Inicializando bot")

	name := "test"
	currency := domain.CurrencyBTC
	targetCurrency := domain.CurrencyUSDT
	initialCapital := decimal.NewFromFloat(1000)
	delta := decimal.NewFromFloat(200)
	monitorInterval := 20 * time.Second

	bot, err := domain.CreateBot(name, currency, targetCurrency, initialCapital, delta, monitorInterval)
	if err != nil {
		return err
	}

	logs.Info(ctx, fmt.Sprintf("Bot %s inicializado, capital_inicial: %s, delta: %s", bot.Name, bot.InitialCapital.String(), bot.Delta))

	go s.executeStrategy(ctx, bot)

	return nil
}

func (s *Init) executeStrategy(ctx context.Context, bot *domain.Bot) {
	for {
		price, err := s.providerRepository.GetPrice(ctx, bot.TargetCurrency, bot.Currency)
		if err != nil {
			logs.Error(ctx, "could not get price", logs.NewAttr("error", err))
		}

		priceRange := bot.CalculatePriceRange(price.Price)

		/**TODO: only for simulate sell with test */
		if bot.HasOpenOrder(priceRange) {
			logs.Info(ctx, fmt.Sprintf("Ya hay una order para el rango de precio de %s", price.Price.String()))
			time.Sleep(bot.MonitorInterval)
			continue
		}

		newOrder, err := bot.GenerateOrder(price.Price, priceRange)
		if err != nil {
			logs.Error(ctx, "could not get generate order", logs.NewAttr("error", err))
		}

		externalId, err := s.providerRepository.CreateOrderInProvider(ctx, newOrder)
		if err != nil {
			logs.Error(ctx, "could not get generate order in provider", logs.NewAttr("error", err))
		}

		newOrder.AddExternalId(externalId)

		// logs.Info(ctx, fmt.Sprintf("order created in provider with external_id: %s", externalId))

		time.Sleep(bot.MonitorInterval)
	}
}
