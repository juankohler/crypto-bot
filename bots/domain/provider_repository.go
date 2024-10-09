package domain

import (
	"context"
)

type ProviderRepository interface {
	GetPrice(ctx context.Context, baseCurrency string, quoteCurrency string) (*Price, error)
	CreateOrderInProvider(ctx context.Context, order *Order, botName string) (string, error)
}
